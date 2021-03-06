package core

import (
	"bufio"
	"context"
	"io"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"time"

	"fmt"

	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/docker/docker/api/types"
	"github.com/gogo/protobuf/proto"
)

// verify all containers to open logs stream if not already done
func (a *Agent) updateLogsStream() {
	for ID, data := range a.containers {
		if data.logsStream == nil || data.logsReadError {
			lastTimeID := a.getLastTimeID(ID)
			if lastTimeID == "" {
				log.Printf("open logs stream from the begining on container %s\n", data.name)
			} else {
				log.Printf("open logs stream from time_id=%s on container %s\n", lastTimeID, data.name)
			}
			stream, err := a.openLogsStream(ID, lastTimeID)
			if err != nil {
				log.Printf("Error opening logs stream on container: %s\n", data.name)
			} else {
				data.logsStream = stream
				go a.startReadingLogs(ID, data)
			}
		}
	}
}

// open a logs container stream
func (a *Agent) openLogsStream(ID string, lastTimeID string) (io.ReadCloser, error) {
	containerLogsOptions := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: true,
	}
	if lastTimeID != "" {
		containerLogsOptions.Since = lastTimeID
	}
	return a.dock.GetClient().ContainerLogs(context.Background(), ID, containerLogsOptions)
}

// get last timestamp if exist
func (a *Agent) getLastTimeID(ID string) string {
	data, err := ioutil.ReadFile(path.Join(containersDataDir, ID))
	if err != nil {
		return ""
	}
	return string(data)
}

// stream reading loop
func (a *Agent) startReadingLogs(ID string, data *ContainerData) {
	stream := data.logsStream
	reader := bufio.NewReader(stream)
	data.lastDateSaveTime = time.Now()
	log.Printf("start reading logs on container: %s\n", data.name)
	var previous, now int64
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				log.Printf("Stream log EOF container terminated: %s\n", data.name)
			} else {
				log.Printf("error reading logs, closing logs stream on container %s (%v)\n", data.name, err)
			}
			data.logsReadError = true
			_ = stream.Close()
			a.removeContainer(ID)
			return
		}
		if len(line) <= 39 {
			// mt.Printf("invalid log: [%s]\n", line)
			continue
		}
		date := line[8:38]
		slog := strings.TrimSuffix(line[39:], "\n")
		timestamp, err := time.Parse("2006-01-02T15:04:05.000000000Z", date)
		if err != nil {
			timestamp = time.Now()
		}
		now = time.Now().UnixNano()
		if now <= previous {
			now = previous + 1
		}
		previous = now
		logEntry := logs.LogEntry{
			Timestamp:          timestamp.Format(time.RFC3339Nano),
			ContainerId:        ID,
			ContainerName:      data.name,
			ContainerShortName: data.shortName,
			ContainerState:     data.state,
			ServiceName:        data.serviceName,
			ServiceId:          data.serviceID,
			TaskId:             data.taskID,
			StackName:          data.stackName,
			NodeId:             data.nodeID,
			TimeId:             fmt.Sprintf("%016X", now),
			Labels:             data.labels,
			Msg:                slog,
		}
		a.addLogEntry(&logEntry, data, date)
	}
}

func (a *Agent) addLogEntry(entry *logs.LogEntry, data *ContainerData, date string) {
	if conf.logsBufferPeriod == 0 || conf.logsBufferSize == 0 {
		a.logsBuffer.Entries[0] = entry
		a.sendLogsBuffer()
		a.periodicDataSave(data, date)
		return
	}
	a.logsBufferMutex.Lock()
	defer a.logsBufferMutex.Unlock()
	if a.logsBuffer == nil {
		a.logsBuffer.Entries = make([]*logs.LogEntry, conf.logsBufferSize)
	}
	a.logsBuffer.Entries = append(a.logsBuffer.Entries, entry)
	if len(a.logsBuffer.Entries) >= conf.logsBufferSize {
		a.sendLogsBuffer()
		a.logsBuffer.Entries = nil
		a.periodicDataSave(data, date)
	}
}

func (a *Agent) sendLogsBuffer() {
	encoded, err := proto.Marshal(a.logsBuffer)
	if err != nil {
		log.Printf("error marshalling log entries: %v\n", err)
		return
	}
	_, err = a.natsStreaming.GetClient().PublishAsync(ns.LogsSubject, encoded, nil)
	if err != nil {
		log.Printf("error sending log entry: %v\n", err)
		return
	}
	a.nbLogs += len(a.logsBuffer.Entries)
}

func (a *Agent) periodicDataSave(data *ContainerData, date string) {
	now := time.Now()
	if now.Sub(data.lastDateSaveTime).Seconds() >= float64(a.logsSavedDatePeriod) {
		err := ioutil.WriteFile(path.Join(containersDataDir, data.ID), []byte(date), 0666)
		if err != nil {
			log.Println("error writing to container data directory: ", err)
		}
		data.lastDateSaveTime = now
	}
}

// close all logs stream
func (a *Agent) closeLogsStreams() {
	for _, data := range a.containers {
		if data.logsStream != nil {
			err := data.logsStream.Close()
			if err != nil {
				log.Println("Error closing a log stream: ", err)
			}
		}
	}
}
