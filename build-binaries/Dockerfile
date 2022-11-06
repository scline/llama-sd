FROM golang:1.13

WORKDIR /app

COPY build.sh build.sh
#RUN go get -d -v github.com/dropbox/llama
#
#RUN go build -v github.com/dropbox/llama/cmd/collector
#RUN go build -v github.com/dropbox/llama/cmd/reflector
#RUN go build -v github.com/dropbox/llama/cmd/scraper

CMD ["bash", "-c", "bash build.sh"]
