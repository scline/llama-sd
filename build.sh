# Simple script to build containers located in this repo, used to pipeline work later down the line

# Build server
version=`cat $PWD/llama-server/version`
docker build $PWD/llama-server -t smcline06/llama-server:$version
docker build $PWD/llama-server -t smcline06/llama-server:latest

docker push smcline06/llama-server:$version
#docker push smcline06/llama-server:latest

# Build client
version=`cat $PWD/llama-client/version`
docker build $PWD/llama-client -t smcline06/llama-client:$version
docker build $PWD/llama-client -t smcline06/llama-client:latest

docker push smcline06/llama-client:$version
#docker push smcline06/llama-client:latest

# Build scraper
version=`cat $PWD/llama-scraper/version`
docker build $PWD/llama-scraper -t smcline06/llama-scraper:$version
docker build $PWD/llama-scraper -t smcline06/llama-scraper:latest

docker push smcline06/llama-scraper:$version
#docker push smcline06/llama-scraper:latest

# Build probe
version=`cat $PWD/llama-probe/version`
docker build $PWD/llama-probe -t smcline06/llama-probe:$version
docker build $PWD/llama-probe -t smcline06/llama-probe:latest

docker push smcline06/llama-probe:$version
#docker push smcline06/llama-probe:latest
