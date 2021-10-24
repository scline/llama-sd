FROM golang:1.13

WORKDIR /go/src/app

COPY entrypoint.sh entrypoint.sh 
COPY register.go register.go 
COPY lamoid-grazer lamoid-grazer

RUN go get -d -v github.com/dropbox/llama
RUN go install -v github.com/dropbox/llama/cmd/collector
RUN go install -v github.com/dropbox/llama/cmd/reflector

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

ENTRYPOINT [ "/go/src/app/lamoid-grazer" ]
