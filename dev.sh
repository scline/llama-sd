#!/bin/bash
# Simple script to build containers located in this repo, used to pipeline work later down the line

# Tag arm if built from RaspberryPi
ARCH=$(uname -m)
echo "ARCH: $ARCH"

case "$ARCH" in
  armv7*)   tag="arm7-"                             ;;
  arm64)    tag="arm64-"                            ;;
  x86_64)   tag=""                                  ;;
  *)        echo "UNKNOWN ARCH, EXITING"; exit      ;;
esac

echo "TAG: $tag"

# Build server
version=`cat $PWD/llama-server/version`
docker build $PWD/llama-server -t llama-server:${tag}${version}-DEV
docker build $PWD/llama-server -t llama-server:${tag}latest-DEV

# Build client
version=`cat $PWD/llama-client/version`
docker build $PWD/llama-client -t llama-client:${tag}${version}-DEV
docker build $PWD/llama-client -t llama-client:${tag}latest-DEV

# Build scraper
version=`cat $PWD/llama-scraper/version`
docker build $PWD/llama-scraper -t llama-scraper:${tag}${version}-DEV
docker build $PWD/llama-scraper -t llama-scraper:${tag}latest-DEV

# Build probe
version=`cat $PWD/llama-probe/version`
docker build $PWD/llama-probe -t llama-probe:${tag}${version}-DEV
docker build $PWD/llama-probe -t llama-probe:${tag}latest-DEV
