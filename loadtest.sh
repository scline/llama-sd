# Loop through device registration, lets see how far we get
i=0
while true 
do
curl -H 'Content-Type: application/json' -X POST \
    -d '{
	"port": '$i',
	"keepalive": 600,
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
