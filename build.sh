# Simple script to build containers located in this repo, used to pipeline work later down the line

# Build server
version=`cat $PWD/server/version`
docker build $PWD/server -t llama-server:$version
docker build $PWD/server -t llama-server:latest

# Build server
version=`cat $PWD/client/version`
docker build $PWD/server -t llama-client:$version
docker build $PWD/server -t llama-client:latest
