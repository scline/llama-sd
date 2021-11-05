# LLAMA Service Discovery
This Python/Docker-based code is adding additional functionality to a network monitoring tool. [LLAMA](https://github.com/dropbox/llama) by itself is fantastic if deploying to a small subset or non-changing environment where IP addresses don't often change, for example.

The added code here allows the deployment of this tool to be a little more dynamic. With a few environment variables, you can, for example, deploy probes to 100's of locations or hosts. They will automatically register and start testing between each other without manual configurations or intervention from the user.

![enter image description here](https://github.com/scline/llama-sd/blob/master/docs/001.gif) 

## What is LLAMA?
LLAMA (Loss and Latency Matrix) is a library for testing and measuring network loss and latency between distributed endpoints.

It does this by sending UDP datagrams/probes from collectors to reflectors and measuring how long it takes for them to return, if they return at all. UDP is used to provide ECMP hashing over multiple paths (a win over ICMP) without the need for setup/teardown and per-packet granularity (a win over TCP).

This was developed and created by DropBox: [Github Project](https://github.com/dropbox/llama)


## Components
### LLAMA-SERVER
The server component is a basic Python3 Flask application serving API endpoints. Its primary function is to accept registration messages in JSON from remote clients and group/present the hosts in a format LLAMA Collectors understand.

#### Script Arguments and Environment Variables
- `-c, --config, APP_CONFIG` - [Configuration file](https://github.com/scline/llama-sd/blob/master/llama-server/src/config.yml) path
- `-g, --group, APP_GROUP` - Default group probes will be assigned if none is given. Probe settings will overwrite this value.
- `-i, --host, APP_HOST` - Server IP to listen for web traffic, 0.0.0.0 is all available IP's. Defaults to 127.0.0.1 if not set.
- `-k, --keepalive, APP_KEEPALIVE`  - Keepalive settings, server will remove probe entries if they do not check in within this window value in seconds. This value is used if probes do not give one.
- `-p, --port, APP_PORT` - Port web server listens on. Defaults to 5000 if not set.
- `-v, --verbose, APP_VERBOSE` - Enable debug logging

### LLAMA- SCRAPER
The Scraper component (written by DropBox) will periodically query a LLAMA-Collector via `http://<probe_ip>:8100/influxdata`. It received JSON metrics about the tests that the probe performed and stored them in InfluxDB

Example of what one of these payloads looks like
```
[
  {
    "fields": {
      "loss": 100.000000,
      "lost": 250.000000,
      "rtt": 0.000000,
      "sent": 250.000000
    },
    "tags": {
      "dst_ip": "104.225.250.92",
      "name": "llama_server",
      "src_ip": "::",
      "src_name": "las2",
      "version": "0.0.1"
    },
    "time": "0001-01-01T00:00:00Z",
    "measurement": "raw_stats"
  }
]
```
#### Environment Variables
- `INFLUXDB_HOST` - The IP or hostname of the influxDB to store metrics, using version 1.8 is recommended.
- `INFLUXDB_NAME` - InfluxDB name where data is stored.
- `INFLUXDB_PORT` - InfluxDB listening port
- `LLAMA_SERVER` - URL of LLAMA Server endpoint for gathering host list. i.e. `http://llama.somehost.com:8081`

#### Groups
You can have multiple groups of probes to one server. Assigning  a group name of `BareMetal` vs `WAN` for example. All nodes in the WAN group will perform a full-mesh test against each other while the `BareMetal` group will do the same for probes registered as such. This allows segmentation and future scaling considerations.
![enter image description here](https://github.com/scline/llama-sd/blob/master/docs/groups.png) 

#### Script Arguments and Environment Variables
- `-c, --config, APP_CONFIG` - [Configuration file](https://github.com/scline/llama-sd/blob/master/llama-client/src/config.yml) path
- `-g, --group, LLAMA_GROUP` - Group the probe will be assigned to.
- `-i, --ip, LLAMA_SOURCE_IP` - Optional, if the client wants to tell the server what the probe IP is. By default the server will grab this information from the API call. This option is required if running servers and clients on the same host (docker IP mess).
- `-k, --keepalive, LLAMA_KEEPALIVE`  - Keepalive settings, server will remove probe entries if they do not check in within this window value in seconds.
- `-s, --server, LLAMA_SERVER` -  URL of LLAMA Server endpoint for gathering host list. i.e. `http://llama.somehost.com:8081`
- `-t, --tags, LLAMA_TAGS` - Additional tags reported by probes. Example: `[ location: Las Vegas, src_name: Office ]`
- `-v, --verbose, APP_VERBOSE` - Enable debug logging

### LLAMA-PROBE
Docker container that contains two LLAMA components created by Dropbox. LLAMA-Reflector and the LLAMA-Collector. Added a small [bash script](https://github.com/scline/llama-sd/blob/master/llama-probe/entrypoint.sh) that will periodically (every 30 seconds) pull a config for whatever Group the probe is part of. If the configuration pulled is different from the running one, the service will restart to consume the changes.

#### Environment Variables
- `LLAMA_SERVER` - URL of LLAMA Server endpoint for gathering host list. i.e. `http://llama.somehost.com:8081`
- `LLAMA_GROUP` - Group name of probe, optional
- `LLAMA_KEEPALIVE` - How long should the cluster ping if probe is unreachable in seconds
- `PROBE_NAME` - Generally a hostname that is tagged on metrics
- `PROVE_SHORTNAME` - Shorter name (i.e. pdx1 for a datacenter in Portland or usw2_1 for an AWS location)

## Installation
Installation via Docker containers is going to be the simplest way. This will work for x86 or ARM-based systems like the Raspberry Pi.

### Copy-Paste Probe install (Linux)
```
docker run --restart unless-stopped -d \
-p 8100:8100/tcp \
-p 8100:8100/udp \
-e LLAMA_SERVER=http://llama.packetpals.com:8105 \
-e LLAMA_GROUP=github \
-e PROBE_NAME=Long_Hostname \
-e PROBE_SHORTNAME=gh1 \
--name llama-probe \
smcline06/llama-probe:latest
```

### Copy-Paste Probe install (Raspberry Pi)
```
docker run --restart unless-stopped -d \
-p 8100:8100/tcp \
-p 8100:8100/udp \
-e LLAMA_SERVER=http://llama.packetpals.com:8105 \
-e LLAMA_GROUP=github \
-e PROBE_NAME=Long_Hostname \
-e PROBE_SHORTNAME=gh1 \
--name llama-probe \
smcline06/llama-probe:arm7-latest
```

## Network Requirements
Probes are hardcoded to use TCP and UDP port 8100 for communication. In the future, this will be configurable. If deploying this behind a NAT, for example, within a SOHO environment, then you will need to set up destination ports accordingly on your home router. 

| Source | Destination | Destination Port | Protocol
|--|--|--|--|
| 0.0.0.0/0 (Internet) | Public IP/Interface |8100 | TCP + UDP| 

![enter image description here](https://github.com/scline/llama-sd/blob/master/docs/network.png) 

## Changelog
n/a (things are very beta)
  
## TODO
- Show examples of using docker-compose
- Outline API URL examples and purposes
- Explain /metrics on Server
