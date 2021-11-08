#!/bin/bash
# Simple bash script to manage llama probe data collection
echo "entrypoint.sh running..."
server_url="$LLAMA_SERVER/api/v1/scraper"

# Output for debugging
echo """
SERVER:         $LLAMA_SERVER
Scraper URL:    $server_url
InfluxDB Host:  $INFLUXDB_HOST
InfluxDB Port:  $INFLUXDB_PORT
InfluxDB Name:  $INFLUXDB_NAME
"""

# Grab a new configuration from the server
# Example output: "1.1.1.1,192.168.1.1,10.0.0.2"
if ! collector_hosts=$(curl -m 5 -s $server_url); then
  # Curl gave an error code
  echo "ERROR: Unable to access '$server_url'"
  exit 1
fi

# Convert string of IP addresses into an array we can iterate against
# "1.1.1.1,2.2.2.2,3.3.3.3" ---> ("1.1.1.1", "2.2.2.2", "3.3.3.3")
collector_array=(`echo $collector_hosts | tr ',' ' '`)
echo "Collector Host List:  ${collector_array[@]}"
echo "Collector Host Count: ${#collector_array[@]}"

# Grab interval metric from API server
interval="`curl -s $LLAMA_SERVER/api/v1/interval`"
echo "Collector Interval: $interval"

echo "Starting Scraper"
for i in "${collector_array[@]}"
do
  # Validate if we are getting an IP address, close if we are not
  if ! [[ "$i" =~ ^(([1-9]?[0-9]|1[0-9][0-9]|2([0-4][0-9]|5[0-5]))\.){3}([1-9]?[0-9]|1[0-9][0-9]|2([0-4][0-9]|5[0-5]))$ ]]; then
    echo "ERROR: '$i' is not a valid IP address. Variable Dump:"
    echo "${collector_array[@]}"
    exit 1
  fi

  # Run scraper thread (one per IP endpoint) 
  # Need to run multiple since it appears the scraper app does not timeout or multi-thread collection. Means one dead probe can stop collection entirely.
  scraper -llama.collector-hosts $i -llama.collector-port 8100 -llama.influxdb-host $INFLUXDB_HOST -llama.influxdb-name $INFLUXDB_NAME -llama.influxdb-port $INFLUXDB_PORT -llama.interval $interval &
done

echo """
########################## Config Checking Loop ##########################

This container will continually loop to see if any changes are available
to the collector list. Since this is a container, we do not care about the
state. The container will close with an exit code of 5 if changes occure. 

Please be sure to have auto-restart enabled!

##########################################################################
"""

while true
do
  # Sleep for 60 seconds
  sleep 60

  # Grab a new configuration from the server
  # Example output: "1.1.1.1,192.168.1.1,10.0.0.2"
  if ! collector_hosts_new=$(curl -m 5 -s $server_url); then
    # Curl gave an error code
    echo "ERROR: Unable to access '$server_url', ignoring...."
    # If we get an error during the config check loop, ignore it and reset. This is to limit an API server error stopping metrics gathering.
    continue
  fi  

  # Compaire host lists
  if [ "$collector_hosts_new" = "$collector_hosts" ]; then
    # Host list did not change
    continue
  else
    # Host list updated!
    echo "Probe list changed!"
    echo "OLD:  $collector_hosts"
    echo "NEW:  $collector_hosts_new"; echo ""
  
    # Stop the container since we should have this set to auto-restart and collect requited changes on the fly.
    echo "Stopping container to pick up new changes."
    exit 5
  fi
done
