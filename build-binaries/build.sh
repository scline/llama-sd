#!/bin/bash
echo "Building golang application"

go get -d -v github.com/dropbox/llama

go build -o /go/bin/collector -v github.com/dropbox/llama/cmd/collector
go build -o /go/bin/reflector -v github.com/dropbox/llama/cmd/reflector
go build -o /go/bin/scraper -v github.com/dropbox/llama/cmd/scraper
