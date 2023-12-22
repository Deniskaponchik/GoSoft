cd Unifi/cmd/unifi

rem swag init -g main.go
rem swag init -g Unifi/internal/controller/http/fokInterface/fokusov.go
swag init


rem go build -o ..\..\bin\Unifi_v3.13 -ldflags="-s -w"
rem go build -o ..\..\bin\Unifi_v3.14-TEST -ldflags="-s -w"
rem go build -o ..\..\bin\Unifi_v3.15-TEST -ldflags="-s -w"
rem go build -o ..\..\bin\Unifi_v3.16-TEST -ldflags="-s -w"
go build -o ..\..\bin\Unifi_v3.17-TEST -ldflags="-s -w"

rem start /B ../../bin/Unifi_v3.13 -mode PROD -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
rem start /B ../../bin/Unifi_v3.14-TEST -mode WEB -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
rem start /B ../../bin/Unifi_v3.15-TEST -mode PROD -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
rem start /B ../../bin/Unifi_v3.16-TEST -mode PROD -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
start /B ../../bin/Unifi_v3.17-TEST -mode PROD -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3

rem start /B MakeD.bat