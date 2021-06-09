# Simple script to build containers located in this repo, used to pipeline work later down the line

# Build server
version=`cat $PWD/llama-server/version`
docker build $PWD/llama-server -t llama-server:$version
docker build $PWD/llama-server -t llama-server:latest

# Build client
version=`cat $PWD/llama-client/version`
docker build $PWD/llama-client -t llama-client:$version
docker build $PWD/llama-client -t llama-client:latest

# Build scraper
version=`cat $PWD/llama-scraper/version`
docker build $PWD/llama-scraper -t llama-scraper:$version
docker build $PWD/llama-scraper -t llama-scraper:latest

# Build probe
version=`cat $PWD/llama-probe/version`
docker build $PWD/llama-probe -t llama-probe:$version
docker build $PWD/llama-probe -t llama-probe:latest
