# Cleanup running things

docker kill llama-server-dev
docker rm llama-server-dev

docker kill llama-client-dev
docker rm llama-client-dev

# Build client
docker build $PWD/client -t llama-client:dev
docker run --name llama-client-dev -d llama-client-dev

# Build server
docker build $PWD/server -t llama-server:dev
docker run --name llama-server-dev -p 80:80 -d llama-server-dev

