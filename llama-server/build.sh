# basic build script for dev

docker kill llama-server-dev
docker rm llama-server-dev

docker build . -t llama-server-dev
docker run --name llama-server-dev -p 80:80 -d llama-server-dev

sleep 3

i=0
while [ $i -lt 3 ] 
do
curl -H 'Content-Type: application/json' -X POST \
    -d '{
	"port": '$i',
	"keepalive": 120,
	"tags": {
		"version": "1.0",
		"src_dc": "SOURCE",
		"dst_datacenter": "Destination"
	}
}' \
    http://127.0.0.1/api/v1/register
((i=i+1))
done

i=0
while [ $i -lt 3 ] 
do
curl -H 'Content-Type: application/json' -X POST \
    -d '{
	"port": '$i',
	"group": "new_group",
	"keepalive": 120,
	"tags": {
		"version": "1.0",
		"src_dc": "SOURCE",
		"dst_datacenter": "Destination"
	}
}' \
    http://127.0.0.1/api/v1/register
((i=i+1))
done