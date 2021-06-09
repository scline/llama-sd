#!/bin/bash
# Simple bash script to manage the llama probe
echo "entrypoint.sh running..."
server_url="$LLAMA_SERVER/api/v1/scraper"

echo "SERVER: $LLAMA_SERVER"
echo "Scraper URL: $server_url"

# Empty variable to start
collector_hosts=""

# Grab a new configuration from the server
echo "Waiting 60 seconds before pulling a list of hosts to scrape..."
sleep 60

# Host list
collector_hosts_new=`curl $server_url`

# Output for debugging
echo "InfluxDB Host: $INFLUXDB_HOST"
echo "InfluxDB Port: $INFLUXDB_PORT"
echo "InfluxDB Name: $INFLUXDB_NAME"
echo "Host List: $collector_hosts_new"

echo "Starting Scraper"
scraper -llama.collector-hosts $collector_hosts_new -llama.collector-port 8100 -llama.influxdb-host $INFLUXDB_HOST -llama.influxdb-name $INFLUXDB_NAME -llama.influxdb-port $INFLUXDB_PORT  -llama.interval 10

echo "~~~ Config Checking Loop ~~~"
while true
do
  # Sleep for 60 seconds, do this first so we dont have issues at startup
  sleep 60

  # Assign old values to variable for later comparing
  collector_hosts=$collector_hosts_new

  # New host list
  collector_hosts_new=`curl $server_url`

  scraper_pid=`ps -A -o pid,cmd | grep scraper | grep -v grep | head -n 1 | awk '{print $1}'`

  if [ -z "$scraper_pid" ]; then
    echo "Scraper process is not running! Restarting..."
    scraper -llama.collector-hosts $collector_hosts_new -llama.collector-port 8100 -llama.influxdb-host $INFLUXDB_HOST -llama.influxdb-name $INFLUXDB_NAME -llama.influxdb-port $INFLUXDB_PORT  -llama.interval 10 &
  fi

  # If (running ISNOTEQUALTO temp)
  if [[ "$collector_hosts" != "$collector_hosts_new" ]]; then
    echo "Config update found!"
    echo "Host List: $collector_hosts_new"

    # Kill running scaper for updated targets
    kill -9 `ps -A -o pid,cmd | grep collector | grep -v grep | head -n 1 | awk '{print $1}'`
    scraper -llama.collector-hosts $collector_hosts_new -llama.collector-port 8100 -llama.influxdb-host $INFLUXDB_HOST -llama.influxdb-name $INFLUXDB_NAME -llama.influxdb-port $INFLUXDB_PORT  -llama.interval 10 &
  fi 

done
