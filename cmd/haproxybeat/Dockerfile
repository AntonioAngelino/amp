FROM ubuntu

RUN apt-get update && \
apt-get install -y curl && \
curl -L -O https://artifacts.elastic.co/downloads/beats/metricbeat/metricbeat-5.4.0-amd64.deb && \
dpkg -i metricbeat-5.4.0-amd64.deb
