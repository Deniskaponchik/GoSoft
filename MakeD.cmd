
rem SET "swag=%~1"
rem if SWAG == "swag"(
rem swag init -g internal/controller/http/fokInterface/routes.go --parseInternal true
rem )

rem if "%2" == "grpc"(
rem D:\Programms\Coding\protoc-25.1-win64\bin\protoc.exe -I api/unifi/v1 api/unifi/v1/unifi.proto ^
rem --go_out=pkg/grpc/unifi/v1 ^
rem --go_opt=paths=source_relative ^
rem --go-grpc_out=pkg/grpc/unifi/v1 ^
rem --go-grpc_opt=paths=source_relative

rem --go_out-plugins-grpc:pkg/grpc/server ^
rem --go_out        сюда складывается результат protoc
rem --go_opt        сгенерированные файлы будут использовать тот же пакет, что и proto-файлы. source_relative означает, что выходные файлы будут иметь тот же пакет, что и исходные .proto файлы
rem --go-grpc_out   куда складывать go grpc код
rem --go-grpc_opt   как создавать имена пакетов для gRPC
rem )

rem if "%3" == "move"(
cd cmd/unifi
rem rem swag init -g ../../internal/controller/http/fokInterface/routes.go -o ../../docs
rem )

rem if "%4" == "build"(
rem go build -o ..\..\bin\Unifi_v3.19-TEST -ldflags="-s -w"
rem go build -o ..\..\bin\Unifi_v3.20-TEST -ldflags="-s -w"
rem go build -o ..\..\bin\Unifi_v3.21 -ldflags="-s -w"
go build -o ..\..\bin\Unifi_v3.22-TEST -ldflags="-s -w"
rem )

rem if "%5" == "run"(
rem start /B ../../bin/Unifi_v3.19-TEST -mode WEB -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
rem start /B ../../bin/Unifi_v3.20-TEST -mode WEB -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
rem start /B ../../bin/Unifi_v3.21 -mode PROD -time +5 -httpUrl 10.57.179.121:8081 -db it_support_db_3
start /B ../../bin/Unifi_v3.22-TEST -mode PROD -time +5
rem )

rem start /B MakeD.bat
rem start /B MakeD nil grpc nil nil nil
rem MakeD