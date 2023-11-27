cd Unifi/cmd/unifi
go build -o ..\..\bin\Unifi_v3.0-TEST -ldflags="-s -w"
start /B ../../bin/Unifi_v3.0-TEST -mode PROD -time +5 -httpUrl 10.57.179.121:8081