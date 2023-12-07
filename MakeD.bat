cd Unifi/cmd/unifi

go build -o ..\..\bin\Unifi_v3.6-TEST -ldflags="-s -w"

rem start /B ../../bin/Unifi_v3.6-TEST -mode PROD -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
start /B ../../bin/Unifi_v3.6-TEST -mode PROD -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3