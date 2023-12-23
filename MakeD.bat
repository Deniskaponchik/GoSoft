rem swag init -g internal/controller/http/fokInterface/routes.go --parseInternal true
cd cmd/unifi
rem swag init -g ../../internal/controller/http/fokInterface/routes.go -o ../../docs

rem go build -o ..\..\bin\Unifi_v3.16-TEST -ldflags="-s -w"
rem go build -o ..\..\bin\Unifi_v3.17-TEST -ldflags="-s -w"
go build -o ..\..\bin\Unifi_v3.18-TEST -ldflags="-s -w"

rem start /B ../../bin/Unifi_v3.16-TEST -mode PROD -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
rem start /B ../../bin/Unifi_v3.17-TEST -mode PROD -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
start /B ../../bin/Unifi_v3.18-TEST -mode WEB -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3

rem start /B MakeD.bat