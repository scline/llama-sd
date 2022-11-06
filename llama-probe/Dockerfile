FROM alpine:3.16

WORKDIR /go/src/app

COPY bin/ /usr/local/bin/
COPY entrypoint.sh entrypoint.sh 
COPY register.go register.go 

# Install dependencies for applications to run
RUN /bin/sh -c "apk add --no-cache gcompat bash curl"

# Install additional network tools in the event troubleshooting is required
RUN /bin/sh -c "apk add --no-cache traceroute bind-tools mtr"

# environment variables
ENV \
	PROBE_NAME="probe_name" \
	PROBE_SHORTNAME="probe_shortname" \
    LLAMA_SERVER="http://127.0.0.1" \
    LLAMA_PORT=8100 \
    LLAMA_KEEPALIVE=300 \
    #LLAMA_SOURCE_IP=10.0.0.1 \
    LLAMA_GROUP="default"

EXPOSE 8100/tcp
EXPOSE 8100/udp

CMD ["bash", "-c", "bash entrypoint.sh"]
