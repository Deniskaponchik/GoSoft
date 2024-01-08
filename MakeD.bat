swag init -g internal/controller/http/fokInterface/routes.go --parseInternal true

rem D:\Programms\Coding\protoc-25.1-win64\bin\protoc.exe -I api/unifi/v1 api/unifi/v1/unifi.proto ^
--go_out=pkg/grpc/server ^
--go_opt=paths=source_relative ^
--go-grpc-out=pkg/grpc/server ^
--go-grpc_opt=paths=source_relative
rem --go_out-plugins-grpc:pkg/grpc/server ^
rem --go_out        сюда складывается результат protoc
rem --go_opt        сгенерированные файлы будут использовать тот же пакет, что и proto-файлы
rem --go-grpc_out    куда складывать go grpc код



cd cmd/unifi
rem rem swag init -g ../../internal/controller/http/fokInterface/routes.go -o ../../docs

rem go build -o ..\..\bin\Unifi_v3.19-TEST -ldflags="-s -w"
rem go build -o ..\..\bin\Unifi_v3.20-TEST -ldflags="-s -w"
go build -o ..\..\bin\Unifi_v3.21-TEST -ldflags="-s -w"

rem start /B ../../bin/Unifi_v3.19-TEST -mode WEB -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
rem start /B ../../bin/Unifi_v3.20-TEST -mode WEB -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
start /B ../../bin/Unifi_v3.21-TEST -mode WEB -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3

rem start /B MakeD.bat