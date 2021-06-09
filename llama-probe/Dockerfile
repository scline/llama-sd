FROM golang:1.13.3

WORKDIR /go/src/app

copy entrypoint.sh entrypoint.sh 

RUN go get -d -v github.com/dropbox/llama
RUN go install -v github.com/dropbox/llama/cmd/collector
RUN go install -v github.com/dropbox/llama/cmd/reflector

# environment variables
ENV \
    LLAMA_SERVER="http://127.0.0.1" \
    LLAMA_GROUP="default" 

EXPOSE 8100

CMD ["bash", "-c", "bash entrypoint.sh"]
