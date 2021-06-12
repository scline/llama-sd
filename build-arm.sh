# Simple script to build containers located in this repo, used to pipeline work later down the line

# Build server
version=`cat $PWD/llama-server/version`
docker build $PWD/llama-server -t smcline06/llama-server-arm:$version
docker build $PWD/llama-server -t smcline06/llama-server-arm:latest

docker push smcline06/llama-server-arm:$version
docker push smcline06/llama-server-arm:latest

# Build client
version=`cat $PWD/llama-client/version`
docker build $PWD/llama-client -t smcline06/llama-client-arm:$version
docker build $PWD/llama-client -t smcline06/llama-client-arm:latest

docker push smcline06/llama-client-arm:$version
docker push smcline06/llama-client-arm:latest

# Build scraper
version=`cat $PWD/llama-scraper/version`
docker build $PWD/llama-scraper -t smcline06/llama-scraper-arm:$version
docker build $PWD/llama-scraper -t smcline06/llama-scraper-arm:latest

docker push smcline06/llama-scraper-arm:$version
docker push smcline06/llama-scraper-arm:latest

# Build probe
version=`cat $PWD/llama-probe/version`
docker build $PWD/llama-probe -t smcline06/llama-probe-arm:$version
docker build $PWD/llama-probe -t smcline06/llama-probe-arm:latest

docker push smcline06/llama-probe-arm:$version
docker push smcline06/llama-probe-arm:latest
