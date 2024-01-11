package my

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/internal/usecase"
	unifiv1 "github.com/deniskaponchik/GoSoft/pkg/grpc/unifi/v1"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	//"github.com/gookit/slog"
	"log/slog"
	"net"
)

//const logFileName = "Unifi_Grpc_"+time.Now().Format("2006-01-02_15.04.05")+".log"

type GrpcServ struct {
	//представляет собой некую пустую имплементацию всех методов gRPC сервиса.
	//Использование этой структуры помогает обеспечить обратную совместимость при изменении .proto файла.
	//Если мы добавим новый метод в наш .proto файл и заново сгенерируем код,
	//но не реализуем этот метод в serverAPI,
	//то благодаря встраиванию UnimplementedAuthServer наш код все равно будет компилироваться,
	//а новый метод просто вернет ошибку "Not implemented".
	unifiv1.UnimplementedGetAnomaliesServer
	gRPCServer *grpc.Server //Tuzov
	port       int
	urest      usecase.UnifiRestIn //interface. НЕ ИСПОЛЬЗОВАТЬ разыменовыватель *
	//logFileName string
	//logger      *log.Logger
	slogger *slog.Logger //Tuzov
}

func New(uuc *usecase.UnifiUseCase, port int, logFileName string) *GrpcServ { //, log *slog.Logger

	slogger := setupLogger(logFileName)

	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			//logging.StartCall, logging.FinishCall,
			//помимо прочего, мы также хотим логировать тело запроса и ответа.
			//Такая опция далеко не всегда уместна, т.к. в теле запросов может находиться секретная информация
			logging.PayloadReceived, logging.PayloadSent,
		),
		// Add any other option (check functions starting with logging.With).
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			//log.Error("Recovered from panic", slog.Any("panic", p))
			slogger.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	//пока всего один интерсептор, и его обернул grpc.ChainUnaryInterceptor —
	//это некий враппер, который принимает в качестве аргументов набор интерсепторов,
	//а когда приходит одиночный запрос (Unary), запускает все эти интерсепторы поочерёдно
	//(об этом говорит слово Chain в названии).
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		//восстановит и обработает панику, если она случится внутри хэндлера.
		//Полезная штука, ведь мы не хотим, чтобы паника в одном запросе уронила нам весь сервис,
		//остановив обработку даже корректных запросов.
		recovery.UnaryServerInterceptor(recoveryOpts...),
		//logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(slogger), loggingOpts...),
	))
	//Помимо одиночных запросов, gRPC умеет работать также с потоковыми (Stream),
	//и для них мы бы использовали grpc.ChainStreamInterceptor

	//TODO: разобраться, нужно ли:
	//authgrpc.Register(gRPCServer, authService)

	return &GrpcServ{
		gRPCServer: gRPCServer,
		urest:      uuc,
		port:       port,
		//logFileName: logFileName,
		slogger: slogger,
	}
}

// MustRun runs gRPC server and panics if any error occurs.
func (g *GrpcServ) MustRun() {
	if err := g.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (g *GrpcServ) Run() error {
	const op = "grpcapp.Run"

	//l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	//g.log.Info("grpc server started", slog.String("addr", l.Addr().String()))
	g.slogger.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err = g.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stops gRPC server.
func (g *GrpcServ) Stop() {
	const op = "grpcapp.Stop"

	//g.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.Int("port", g.port))
	g.slogger.With(slog.String("op", op)).Info("stopping gRPC server", slog.Int("port", g.port))

	g.gRPCServer.GracefulStop()
}
