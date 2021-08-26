FROM golang:1.13.3

WORKDIR /go/src/app

copy entrypoint.sh entrypoint.sh 

RUN go get -d -v github.com/dropbox/llama
RUN go install -v github.com/dropbox/llama/cmd/scraper

# environment variables
ENV \
    INFLUXDB_HOST="127.0.0.1" \
    INFLUXDB_NAME="llama" \
    INFLUXDB_PORT=8086 \
    LLAMA_INTERVAL=10 \
    LLAMA_PORT=8100 \
    LLAMA_SERVER="http://127.0.0.1"

EXPOSE 8100

CMD ["bash", "-c", "bash entrypoint.sh"]
