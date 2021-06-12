#!/bin/bash
# Simple script to build containers located in this repo, used to pipeline work later down the line

# Tag arm if built from RaspberryPi
OS=$(uname -m)

case "$OS" in
  arm7*)   tag="arm7" ;;
  arm64*)   tag="arm64" ;;
  *)      tag="" ;;
esac

echo "TAG: $tag"

# Build server
version=`cat $PWD/llama-server/version`
docker build $PWD/llama-server -t smcline06/llama-server:${tag}${version}
docker build $PWD/llama-server -t smcline06/llama-server:${tag}latest

docker push smcline06/llama-server:$tag_$version
docker push smcline06/llama-server:${tag}latest

# Build client
version=`cat $PWD/llama-client/version`
docker build $PWD/llama-client -t smcline06/llama-client:$tag_$version
docker build $PWD/llama-client -t smcline06/llama-client:${tag}latest

docker push smcline06/llama-client:${tag}${version}
docker push smcline06/llama-client:${tag}latest

# Build scraper
version=`cat $PWD/llama-scraper/version`
docker build $PWD/llama-scraper -t smcline06/llama-scraper${tag}${version}
docker build $PWD/llama-scraper -t smcline06/llama-scraper:${tag}latest

docker push smcline06/llama-scraper${tag}${version}
docker push smcline06/llama-scraper:${tag}latest

# Build probe
version=`cat $PWD/llama-probe/version`
docker build $PWD/llama-probe -t smcline06/llama-probe${tag}${version}
docker build $PWD/llama-probe -t smcline06/llama-probe:${tag}latest

docker push smcline06/llama-probe${tag}${version}
docker push smcline06/llama-probe:${tag}latest
