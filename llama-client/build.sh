# basic build script for dev

docker kill llama-client-dev
docker rm llama-client-dev

docker build . -t llama-client-dev
docker run --name llama-client-dev llama-client-dev
