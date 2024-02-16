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
rem go build -o ..\..\bin\Unifi_v3.25-TEST -ldflags="-s -w"
rem go build -o ..\..\bin\Unifi_v3.26-TEST -ldflags="-s -w"
rem go build -o ..\..\bin\Unifi_v3.27-PROD -ldflags="-s -w"
    go build -o ..\..\bin\Unifi_v3.28-TEST -ldflags="-s -w"
rem )

cd ../..
rem cd D:\Clouds\GitHub\GoSoft\bin

rem if "%5" == "run"(
rem start /B bin/Unifi_v3.25-TEST -mode PROD -time +5
rem start /B bin/Unifi_v3.26-TEST -mode PROD -time +5
rem start /B bin/Unifi_v3.27-PROD -mode PROD -time +5
    start /B bin/Unifi_v3.28-TEST -mode PROD -time +5
rem )

rem start /B ./bin/Unifi_v3.23-TEST -mode PROD -time +5
rem start /B Unifi_v3.24-TEST -mode PROD -time +5
rem start /B MakeD.bat
rem start /B MakeD nil grpc nil nil nil
rem MakeD