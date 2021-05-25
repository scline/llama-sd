# Loop through device registration, lets see how far we get
i=0
while true 
do
curl -H 'Content-Type: application/json' -X POST \
    -d '{
	"id": "'"$i"'",
	"address": "192.168.0.1",
	"port": 8100,
	"keepalive": 120,
	"group": "testing",
	"meta": {
		"version": "1.1",
		"src_dc": "SOURCE",
		"dst_datacenter": "Destination"
	}
}' \
    http://127.0.0.1/api/v1/register
((i=i+1))
done
