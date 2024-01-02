swag init -g internal/controller/http/fokInterface/routes.go --parseInternal true
cd cmd/unifi
rem rem swag init -g ../../internal/controller/http/fokInterface/routes.go -o ../../docs

rem go build -o ..\..\bin\Unifi_v3.19-TEST -ldflags="-s -w"
rem go build -o ..\..\bin\Unifi_v3.20-TEST -ldflags="-s -w"
go build -o ..\..\bin\Unifi_v3.21-TEST -ldflags="-s -w"

rem start /B ../../bin/Unifi_v3.19-TEST -mode WEB -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
rem start /B ../../bin/Unifi_v3.20-TEST -mode WEB -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
start /B ../../bin/Unifi_v3.21-TEST -mode WEB -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3

rem start /B MakeD.bat