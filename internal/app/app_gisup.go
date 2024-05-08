package app

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/internal/usecase/netdial"

	//"github.com/deniskaponchik/GoSoft/config/ui"
	"github.com/deniskaponchik/GoSoft/config/gisup"
	myGRPC "github.com/deniskaponchik/GoSoft/internal/controller/grpc/my"
	fokInterface "github.com/deniskaponchik/GoSoft/internal/controller/http/fokInterface"
	"github.com/deniskaponchik/GoSoft/internal/usecase"
	"github.com/deniskaponchik/GoSoft/internal/usecase/amqp_rmq"
	"github.com/deniskaponchik/GoSoft/internal/usecase/api_rest"
	"github.com/deniskaponchik/GoSoft/internal/usecase/api_soap"
	"github.com/deniskaponchik/GoSoft/internal/usecase/api_web"
	"github.com/deniskaponchik/GoSoft/internal/usecase/authentication"
	"github.com/deniskaponchik/GoSoft/internal/usecase/authorization"
	"github.com/deniskaponchik/GoSoft/internal/usecase/repo"
	"github.com/deniskaponchik/GoSoft/pkg/postgres"
	"os"
	"os/signal"
	"syscall"

	//"github.com/deniskaponchik/GoSoft/pkg/logger"
	//"github.com/deniskaponchik/GoSoft/pkg/rabbitmq/rmq_rpc/server"
	"log"
	"time"
)

// Run creates objects via constructors.
// func RunUnifi(cfg *ui.ConfigUi) {
func RunGisup(cfg *gisup.ConfigGisup) {

	//Postgres
	pg, err := postgres.New(cfg.PG.PgConnectStr)
	if err != nil {
		log.Printf("Подключение к Postgres завершилось ошибкой: %w", err)
	} else {
		log.Println("Проверка подключения к БД Postgres прошла успешно")
	}
	defer pg.Close()

	//MySQL
	gisupRepo, err := repo.NewGisupRepo(cfg.GLPI.GlpiConnectStr, cfg.GLPI.DB)
	if err != nil {
		log.Fatalf("Подключение к GLPI завершилось ошибкой: %w", err)
	} else {
		log.Println("Проверка подключения к БД GLPI прошла успешно")
	}

	//SOAP
	gisupSoap, err := api_soap.NewSoap(cfg.SoapUrl, cfg.BpmUrl) // cfg.SoapTest, cfg.BpmTest
	if err != nil {
		log.Fatalf("Подключение к SOAP завершилось ошибкой: %w", err)
	} else {
		log.Println("Проверка подключения к SOAP прошла успешно")
	}

	//LDAP
	gisupLdap, err := authentication.NewLdap(cfg.LdapDN, cfg.LdapDomain, cfg.LdapLogin, cfg.LdapPassword, cfg.LdapRoleDn, cfg.LdapServer),
	if err != nil {
		log.Println("Подключение к LDAP завершилось ошибкой: %w", err)
	} else {
		log.Println("Проверка подключения к SOAP прошла успешно")
	}

	//JWT
	gisupJwt, err := authorization.NewAuthJwt(cfg.Token.JwtKey, cfg.Token.TTL)
	if err != nil {
		log.Println("Создание JWT-токена завершилось ошибкой: %w", err)
	} else {
		log.Println("Создание JWT-токена прошло успешно")
	}

	//C3PO
	gisupC3po := api_rest.NewGisupC3po(cfg.C3po.C3poUrl)
	err = gisupC3po.GetUserLogin()
	if err != nil {
		log.Println("Подключение к C3PO завершилось ошибкой: %w", err)
	} else {
		log.Println("Проверка подключения к C3PO прошла успешно")
	}

	//RMQ
	gisupRmq := amqp_rmq.NewRmqUnifi(cfg.RMQ.RmqConnectStr, cfg.RMQ.ServerExchange)
	err = gisupRmq.Publish("Start", "")
	if err != nil {
		log.Println("Подключение к RMQ завершилось ошибкой: %w", err)
	} else {
		log.Println("Проверка подключения к RMQ прошла успешно")
	}


	log.Println("")

	//USECASE
	unifiUseCase := usecase.NewUnifiUC(
		unifiRepo, //вставляем объект, который удовлетворяет интерфейсу UnifiRepo
		amqp_rmq.NewRmqUnifi(cfg.RMQ.RmqConnectStr, cfg.RMQ.ServerExchange),
		api_soap.NewSoap(cfg.SoapUrl, cfg.BpmUrl), // cfg.SoapTest, cfg.BpmTest
		api_web.NewUi(cfg.Ubiquiti.UiUsername, cfg.Ubiquiti.UiPassword, cfg.Ubiquiti.UiContrlRostov, 1),
		api_web.NewUi(cfg.Ubiquiti.UiUsername, cfg.Ubiquiti.UiPassword, cfg.Ubiquiti.UiContrlNovosib, 2),
		api_rest.NewUnifiC3po(cfg.C3po.C3poUrl),
		authentication.NewLdap(cfg.LdapDN, cfg.LdapDomain, cfg.LdapLogin, cfg.LdapPassword, cfg.LdapRoleDn, cfg.LdapServer),
		authorization.NewAuthJwt(cfg.Token.JwtKey, cfg.Token.TTL),

		cfg.App.EveryCodeMap,
		cfg.App.TimeZone,
		cfg.HTTP.URL,

		cfg.Ubiquiti.Daily,
		cfg.Ubiquiti.H1,
		cfg.Ubiquiti.H2,
	)
	go unifiUseCase.InfinityProcessingUnifi()
	/*https://stackoverflow.com/questions/25142016/how-to-return-a-error-from-a-goroutine-through-channels
	errors := make(chan error, 0)
	go func() {
		err = unifiUseCase.InfinityProcessingUnifi()
		if err != nil {
			errors <- err
			return
		}
	}()
	*/

	eltexUseCase := usecase.NewEltex(
		)

	polyRepo, err := repo.NewPolyRepo(cfg.GLPI.GlpiConnectStr, cfg.GLPI.DB)
	if err != nil {
		//если БД недоступна - останавливаем тут же
		l.Fatal(fmt.Errorf("app - Run - glpi.New: %w", err))
	} else {
		fmt.Println("Проверка подключения к БД прошла успешно")
	}

	polyUseCase := usecase.NewPoly(
		//repo.New(cfg.GLPI.GlpiConnectStrITsupport),
		polyRepo,
		//webapi.New(cfg.PolyUsername, cfg.PolyPassword),
		//api_web.NewPolyWebApi(cfg.PolyUsername, cfg.PolyPassword),
		api_rest.NewPolyWebApi(cfg.PolyUsername, cfg.PolyPassword),
		netdial.New(),
		api_soap.New(cfg.SoapUrl, cfg.BpmUrl), // cfg.SoapTest, cfg.BpmTest
		cfg.InnerVars.EveryCodeMap,
		cfg.InnerVars.RestartHour,
		cfg.App.TimeZone,
	)

	gisupUseCase := usecase.NewGisup(
		wifiUseCase,
		vcsUseCase,
		)

	//GRPC
	myGrpc := myGRPC.New(
		//unifiUseCase,
		gisupUseCase,
		cfg.GRPC.Port,
		"logs/Unifi_Grpc_"+time.Now().Format("2006-01-02_15.04.05")+".log",
	)
	go myGrpc.MustRun()
	//go application.GRPCServer.MustRun()
	/*olezhek
	olezhekGRPC := olezhekClean.New(
		unifiUseCase,
		cfg.GRPC.Port,
		"Unifi_Grpc_"+time.Now().Format("2006-01-02_15.04.05")+".log",
	)
	olezhekGRPC.Start()
	*/

	// RabbitMQ RPC Server
	/*EVRONE
	rmqRouter := amqprpc.NewRouter(unifiUseCase)
	//rmqServer, err := server.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
	rmqServer, err := rmqrpcserv.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	}*/

	//HTTP v1
	//FOKUSOV
	httpFokusov := fokInterface.New(
		//gin.Engine,
		unifiUseCase,
		//usecase.Rest(),
		cfg.HTTP.Port,
		//"Unifi_Gin_"+time.Now().Format("2006-01-02_15.04.05")+".log",
		"logs/Unifi_Gin_"+time.Now().Format("2006-01-02_15.04.05")+".log",
		cfg.Token.TTL, //синхронизация времени жизни куки и токена
	)
	httpFokusov.Start()
	/* EVRONE
	handler := gin.New()
	v1.NewRouter(handler, l, unifiUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	*/

	// Graceful shutdown
	//EVRONE
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	s := <-interrupt
	log.Println("Получен сигнал с клавиатуры на остановку работы приложения" + s.String())

	httpFokusov.Stop()
	myGrpc.Stop()
	/*
		err = rmqServer.Shutdown()
		if err != nil {
			l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
		}*/

	/* Tuzov
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	application.GRPCServer.Stop()
	log.Info("Gracefully stopped")
	*/
}
