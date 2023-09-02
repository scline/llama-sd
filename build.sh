#!/bin/bash
# Simple script to build containers located in this repo, used to pipeline work later down the line

# Build server
version=`cat $PWD/llama-server/version`
docker build $PWD/llama-server -t smcline06/llama-server:${tag}${version}
docker build $PWD/llama-server -t smcline06/llama-server:${tag}latest

#docker push smcline06/llama-server:${tag}${version}
#docker push smcline06/llama-server:${tag}latest

# Build scraper
version=`cat $PWD/llama-scraper/version`
docker build $PWD/llama-scraper -t smcline06/llama-scraper:${tag}${version}
docker build $PWD/llama-scraper -t smcline06/llama-scraper:${tag}latest

#docker push smcline06/llama-scraper:${tag}${version}
#docker push smcline06/llama-scraper:${tag}latest

# Build probe
version=`cat $PWD/llama-probe/version`
docker build $PWD/llama-probe -t smcline06/llama-probe:${tag}${version}
docker build $PWD/llama-probe -t smcline06/llama-probe:${tag}latest

#docker push smcline06/llama-probe:${tag}${version}
#docker push smcline06/llama-probe:${tag}latest
