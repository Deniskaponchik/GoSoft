cd Unifi/cmd/unifi
rem go build -o ..\..\bin\Unifi_v3.11-TEST -ldflags="-s -w"
rem go build -o ..\..\bin\Unifi_v3.12-TEST -ldflags="-s -w"
go build -o ..\..\bin\Unifi_v3.13-TEST -ldflags="-s -w"
rem start /B ../../bin/Unifi_v3.11-TEST -mode PROD -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
rem start /B ../../bin/Unifi_v3.12-TEST -mode PROD -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
start /B ../../bin/Unifi_v3.13-TEST -mode PROD -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3

rem start /B MakeD.bat