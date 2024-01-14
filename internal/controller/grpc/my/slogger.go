package my

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"io"
	"log"
	//"github.com/gookit/slog"
	//"github.com/gookit/slog/handler"
	"log/slog"
	"os"
)

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...) //go version 1.21
		//l.Log(slog.Level(lvl), fields...)			 //go version early 1.21
	})
}

func setupLogger(logFileName string) *slog.Logger {

	//fileLogGrpc, err := os.OpenFile(g.logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	fileLogGrpc, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	multiWriter := io.MultiWriter(os.Stdout, fileLogGrpc)

	//var log *slog.Logger
	var slogger *slog.Logger

	//log = slog.New(slog.NewTextHandler(os.Stderr, nil))
	//log = slog.New(slog.NewTextHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slogger = slog.New(slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelDebug}))
	//log = slog.New(slog.NewSimpleHandler(os.Stderr, nil))

	/*Tuzov
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}*/

	return slogger //log
}

/*
func getMultiwriter() *io {}

// for github.com/gookit/slog not for log/slog
func method1() *slog.Logger {

	//h1 := handler.NewConsoleHandler(slog.AllLevels)
	h1 := handler.NewSimpleHandler(slog.AllLevels)

	l := slog.New(h1)
	l.AddHandlers(h1)
}*/
