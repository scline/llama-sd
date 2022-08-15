#FROM golang:1.13

FROM alpine:3.16

WORKDIR /go/src/app

COPY bin/ /usr/local/bin/
COPY entrypoint.sh entrypoint.sh 

# Install dependencies for applications to run
RUN /bin/sh -c "apk add --no-cache gcompat bash curl"

# Install additional network tools in the event troubleshooting is required
RUN /bin/sh -c "apk add --no-cache traceroute bind-tools mtr"

# environment variables
ENV \
    INFLUXDB_HOST="127.0.0.1" \
    INFLUXDB_NAME="llama" \
    INFLUXDB_PORT=8086 \
    LLAMA_SERVER="http://127.0.0.1"

EXPOSE 8100

CMD ["bash", "-c", "bash entrypoint.sh"]
