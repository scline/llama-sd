#!/bin/bash
# Simple bash script to manage the llama probe
echo "entrypoint.sh running..."
server_url="$LLAMA_SERVER/api/v1/config/$LLAMA_GROUP?llamaport=$LLAMA_PORT"

# Stitch source IP variable to URL if provided
if [ "$LLAMA_SOURCE_IP" ]; then
  echo "Probe wants to report its own IP as $LLAMA_SOURCE_IP"
  server_url="$server_url&srcip=$LLAMA_SOURCE_IP"
fi

echo "SERVER: $LLAMA_SERVER"
echo "GROUP: $LLAMA_GROUP"
echo "PORT: $LLAMA_PORT"
echo "KEEPALIVE: $LLAMA_KEEPALIVE"
echo "PROBE NAME: $PROBE_NAME"
echo "PROBE SHORTNAME: $PROBE_SHORTNAME"
echo "Config URL: $server_url"

echo "Starting Reflector"
reflector -port 8100 &

# Run registration GoLang script
echo "Register Probe"
go run register.go

# Grab a new configuration from the server
echo "Waiting 10 seconds before pulling a config file..."
sleep 10

# Save new configuration
curl -s $server_url --output config.yaml

# Output config for docker logging
echo "Configuration file:"
cat config.yaml

echo "Starting Collector"
collector -llama.config config.yaml &

echo "~~~ Config Checking Loop ~~~"
while true
do
  # Run registration GoLang script
  go run register.go

  # Sleep for 30 seconds, do this first so we dont have issues at startup
  sleep 30

  # Grab new config
  curl -s $server_url --output config.yaml.tmp

  # Store MD5 hash of running and canidate config for later validations
  running_config="`md5sum config.yaml | awk '{print $1}'`"
  temp_config="`md5sum config.yaml.tmp | awk '{print $1}'`"

  collector_pid=`ps -A -o pid,cmd | grep collector | grep -v grep | head -n 1 | awk '{print $1}'`

  # If (running ISNOTEQUALTO temp)
  if [[ "$running_config" != "$temp_config" ]]; then
    echo "Config update found!"
    cp -fr config.yaml.tmp config.yaml

    # Send sigup to collector process in order to reload configuration
    # https://github.com/dropbox/llama/blob/master/cmd/collector/main.go#L34
    kill -HUP `ps -A -o pid,cmd | grep collector | grep -v grep | head -n 1 | awk '{print $1}'`
  fi 

  if [ -z "$collector_pid" ]; then
    echo "Collector process is not running! Restarting..."
    collector -llama.config config.yaml &
  fi

  # remove temp files
  rm config.yaml.tmp
done
