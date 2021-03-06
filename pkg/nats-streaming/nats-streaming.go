package ns

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
)

const (
	DefaultURL     = "nats://nats:4222"
	ClusterID      = "test-cluster"
	LogsSubject    = "amp-logs"
	MetricsSubject = "amp-metrics"
)

// NatsStreaming NATS-Streaming wrapper
type NatsStreaming struct {
	client    stan.Conn
	url       string
	clusterID string
	clientID  string
	timeout   time.Duration
	connected bool
}

// NewClient instantiates a NatsStreaming wrapper
func NewClient(url string, clusterID string, clientID string, timeout time.Duration) *NatsStreaming {
	return &NatsStreaming{
		url:       url,
		clusterID: clusterID,
		clientID:  clientID,
		timeout:   timeout,
	}
}

// Connect to NATS-Streaming
func (ns *NatsStreaming) Connect() error {
	if ns.connected {
		return nil
	}
	log.Printf("Connecting to nats streaming, url: %s, clusterId: %s, clientId: %s, timeout: %s\n", ns.url, ns.clusterID, ns.clientID, ns.timeout)
	nc, err := nats.Connect(ns.url, nats.Timeout(ns.timeout))
	if err != nil {
		ns.connected = false
		return fmt.Errorf("unable to connect to nats streaming: %v", err)
	}
	ns.client, err = stan.Connect(ns.clusterID, ns.clientID, stan.NatsConn(nc), stan.ConnectWait(ns.timeout))
	if err != nil {
		ns.connected = false
		return fmt.Errorf("unable to connect to nats streaming: %v", err)
	}
	ns.connected = true
	log.Println("Connected to nats streaming successfuly")
	return nil
}

// GetClient return client
func (ns *NatsStreaming) GetClient() stan.Conn {
	return ns.client
}

// Close the client
func (ns *NatsStreaming) Close() error {
	return ns.client.Close()
}
