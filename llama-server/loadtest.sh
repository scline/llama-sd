# Loop through device registration, lets see how far we get
i=0
while true 
do
curl -H 'Content-Type: application/json' -X POST \
    -d '{
	"port": '$i',
	"keepalive": 900,
	"group": "'$i'",
	"tags": {
		"version": "1.1",
		"src_dc": "SOURCE",
		"dst_datacenter": "Destination"
	}
}' \
    http://10.1.0.107/api/v1/register
((i=i+1))
done
