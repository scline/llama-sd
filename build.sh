# basic build script for dev

docker kill llama-server-dev
docker rm llama-server-dev

docker build . -t llama-server-dev
docker run --name llama-server-dev -p 80:80 -d llama-server-dev

sleep 3

i=0
while [ $i -lt 10 ] 
do
curl -H 'Content-Type: application/json' -X POST \
    -d '{
	"id": "'"$i"'",
	"address": "192.168.0.1",
	"port": 8100,
	"keepalive": 120,
	"group": "testing",
	"meta": {
		"version": "1.0",
		"src_dc": "SOURCE",
		"dst_datacenter": "Destination"
	}
}' \
    http://127.0.0.1/api/v1/register
((i=i+1))
done
