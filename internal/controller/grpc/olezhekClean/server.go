package olezhekClean

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/internal/usecase"
	desc "github.com/deniskaponchik/GoSoft/pkg/grpc/unifi/v1"
	grpcv1 "github.com/deniskaponchik/GoSoft/pkg/grpc/unifi/v1"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"os"
)

type GrpcServ struct {
	//Она представляет собой некую пустую имплементацию всех методов gRPC сервиса.
	//Использование этой структуры помогает обеспечить обратную совместимость при изменении .proto файла.
	//Если мы добавим новый метод в наш .proto файл и заново сгенерируем код,
	//но не реализуем этот метод в serverAPI,
	//то благодаря встраиванию UnimplementedAuthServer наш код все равно будет компилироваться,
	//а новый метод просто вернет ошибку "Not implemented".
	grpcv1.UnimplementedGetAnomaliesServer
	grpcServer *grpc.Server
	urest      usecase.UnifiRestIn //interface. НЕ ИСПОЛЬЗОВАТЬ разыменовыватель *

	port        int
	logFileName string
	logger      *log.Logger
}

func New(uuc *usecase.UnifiUseCase, port int, logFileName string) *GrpcServ {
	return &GrpcServ{
		urest:       uuc,
		port:        port,
		logFileName: logFileName,
	}
}

func (g *GrpcServ) Start() {

	//LOGGER
	//FileNameGin := "Unifi_Gin_" + time.Now().Format("2006-01-02_15.04.05") + ".log"
	fileLogGrpc, err := os.OpenFile(g.logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	multiWriter := io.MultiWriter(os.Stdout, fileLogGrpc)
	gin.DefaultWriter = multiWriter
	//g.grpcServer.
	gin.DefaultErrorWriter = multiWriter
	g.logger = log.New(multiWriter, "", 0)

	//INIT GRPC SERVER
	//a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	g.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	//reflection.Register(a.grpcServer)
	reflection.Register(g.grpcServer)

	//desc.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImpl())
	desc.RegisterGetAnomaliesServer(g.grpcServer, g)

	//RUN SERVER
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		g.logger.Fatalf("Failed to listen/open port: #{err}")
		//g.logger.Fatal("Failed to listen/open port")
	}

	//err = a.grpcServer.Serve(list)
	err = g.grpcServer.Serve(listener)
	if err != nil {
		//return err
		g.logger.Fatalf("Failed to Serve grpc: #{err}")
	}
}

func GetClient() {}

// TODO: GRPC Gracefull shutdown
func (g *GrpcServ) Stop() {}
