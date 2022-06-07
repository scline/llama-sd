#!/bin/bash
version="1.0.4"

# Custom environment settings
if [[ "$MESOS" ]]; then
  echo "MESOS Option Selected!"
  export PROBE_NAME=$MESOS_CONTAINER_IP
  export PROBE_SHORTNAME=$HOST
fi

if [[ "$KUB" ]]; then
  echo "KUB Option Selected!"
  export PROBE_NAME=$HOSTNAME
  export PROBE_SHORTNAME=$RC_HOSTNAME
fi

# Simple bash script to manage the llama probe
echo "entrypoint.sh running..."
server_url="$LLAMA_SERVER/api/v1/config/$LLAMA_GROUP?llamaport=$LLAMA_PORT"
registration_url="$LLAMA_SERVER/api/v1/register"

# Stitch source IP variable to URL if provided
# https://stackoverflow.com/questions/3601515/how-to-check-if-a-variable-is-set-in-bash/16753536
if [ -z ${LLAMA_SOURCE_IP+x} ]; then
  echo "Source IP is not set, will be set by the server on registration."
else
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
echo "Registration URL: $registration_url"
echo "Registration Version: $version"

echo "Starting Reflector"
reflector -port 8100 &

# Run registration 
if [ -z ${LLAMA_SOURCE_IP+x} ]; then
  echo "Register Probe"
  curl -s -X POST $registration_url -H 'Content-Type: application/json' -d "$(cat <<EOF
{ "port": $LLAMA_PORT, "keepalive": $LLAMA_KEEPALIVE, "tags": { "version": "$version", "probe_shortname": "$PROBE_SHORTNAME", "probe_name": "$PROBE_NAME" }, "group": "$LLAMA_GROUP" } 
EOF
)" > /dev/null

  echo "Registration Payload: $(cat <<EOF
{ "port": $LLAMA_PORT, "keepalive": $LLAMA_KEEPALIVE, "tags": { "version": "$version", "probe_shortname": "$PROBE_SHORTNAME", "probe_name": "$PROBE_NAME" }, "group": "$LLAMA_GROUP" } 
EOF
)"

else
  echo "Register Probe w/source IP"
  curl -s -X POST $registration_url -H 'Content-Type: application/json' -d "$(cat <<EOF
{ "port": $LLAMA_PORT, "keepalive": $LLAMA_KEEPALIVE, "ip": $LLAMA_SOURCE_IP, "tags": { "version": "$version", "probe_shortname": "$PROBE_SHORTNAME", "probe_name": "$PROBE_NAME" }, "group": "$LLAMA_GROUP" } 
EOF
)" > /dev/null

  echo "Registration Payload: $(cat <<EOF
{ "port": $LLAMA_PORT, "keepalive": $LLAMA_KEEPALIVE, "ip": $LLAMA_SOURCE_IP, "tags": { "version": "$version", "probe_shortname": "$PROBE_SHORTNAME", "probe_name": "$PROBE_NAME" }, "group": "$LLAMA_GROUP" } 
EOF
)"
fi

# Registration golang script spikes CPU enough to affact latancy on low-CPU environments.
#go run register.go

# Grab a new configuration from the server
echo "Waiting 10 seconds before pulling a config file..."
sleep 10

# Save new configuration
curl -s $server_url --output config.yaml

# Output config for docker logging
echo "Configuration file:"
cat config.yaml

# Store interval value, we need to kill -9 the collector if this changes
interval=`cat config.yaml | grep "interval:" | awk '{print $2}'`

echo "Starting Collector"
collector -llama.config config.yaml &

echo "~~~ Config Checking Loop ~~~"
while true
do
  if [ -z ${LLAMA_SOURCE_IP+x} ]; then
    # Register, no source IP
    curl -s -X POST $registration_url -H 'Content-Type: application/json' -d "$(cat <<EOF
{ "port": $LLAMA_PORT, "keepalive": $LLAMA_KEEPALIVE, "tags": { "version": "$version", "probe_shortname": "$PROBE_SHORTNAME", "probe_name": "$PROBE_NAME" }, "group": "$LLAMA_GROUP" } 
EOF
)" > /dev/null
  else
    # Register with source IP
    curl -s -X POST $registration_url -H 'Content-Type: application/json' -d "$(cat <<EOF
{ "port": $LLAMA_PORT, "keepalive": $LLAMA_KEEPALIVE, "ip": $LLAMA_SOURCE_IP, "tags": { "version": "$version", "probe_shortname": "$PROBE_SHORTNAME", "probe_name": "$PROBE_NAME" }, "group": "$LLAMA_GROUP" } 
EOF
)" > /dev/null
fi

  # Registration golang script spikes CPU enough to affact latancy on low-CPU environments.
  #go run register.go

  # Sleep for 30 seconds
  sleep 30

  # Grab new config
  curl -s $server_url --output config.yaml.tmp

  # Grab new interval
  interval_new=`cat config.yaml.tmp | grep "interval:" | awk '{print $2}'`

  # Store MD5 hash of running and canidate config for later validations
  running_config="`md5sum config.yaml | awk '{print $1}'`"
  temp_config="`md5sum config.yaml.tmp | awk '{print $1}'`"

  collector_pid=`ps -A -o pid,cmd | grep collector | grep -v grep | head -n 1 | awk '{print $1}'`

  # If (running ISNOTEQUALTO temp)
  if [[ "$running_config" != "$temp_config" ]]; then
    echo "Config update found!"
    cp -fr config.yaml.tmp config.yaml

    # If inverval changes kill the collector outright
    if [[ "$interval_new" != "$interval" ]]; then
      echo "Interval has changed, hard-stopping the Collector"

      # Kill -9 the collector since -HUP does not restart with new interval values, then restart
      kill -9 `ps -A -o pid,cmd | grep collector | grep -v grep | head -n 1 | awk '{print $1}'`
      #sleep 5
      #collector -llama.config config.yaml &

      # Set the new interval as base for future runs
      interval=$interval_new

      echo "New Interval: $interval_new"
    fi
    # Send sigup to collector process in order to reload configuration
    # https://github.com/dropbox/llama/blob/master/cmd/collector/main.go#L34
    kill -HUP `ps -A -o pid,cmd | grep collector | grep -v grep | head -n 1 | awk '{print $1}'` 2>/dev/null
  fi 

  if [ -z "$collector_pid" ]; then
    echo "Collector process is not running! Restarting..."
    collector -llama.config config.yaml &
  fi

  # remove temp files
  rm config.yaml.tmp
done
