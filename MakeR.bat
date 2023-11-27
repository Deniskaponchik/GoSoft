cd Unifi/cmd/unifi
go build -o ..\..\bin\Unifi_v3.0-TEST -ldflags="-s -w"
start /B ../../bin/Unifi_v3.0-TEST -mode PROD -cntrl Rostov -time +5 -httpUrl wsir-it-03:8081