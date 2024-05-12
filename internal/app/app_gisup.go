package app

import (
	"github.com/deniskaponchik/GoSoft/internal/usecase/amqp_rmq"
	"github.com/deniskaponchik/GoSoft/internal/usecase/api_web"
	"github.com/deniskaponchik/GoSoft/internal/usecase/authorization"
	"github.com/deniskaponchik/GoSoft/internal/usecase/netdial"
	//"github.com/deniskaponchik/GoSoft/config/ui"
	"github.com/deniskaponchik/GoSoft/config/gisup"
	myGRPC "github.com/deniskaponchik/GoSoft/internal/controller/grpc/my"
	fokInterface "github.com/deniskaponchik/GoSoft/internal/controller/http/fokInterface"
	"github.com/deniskaponchik/GoSoft/internal/usecase"
	"github.com/deniskaponchik/GoSoft/internal/usecase/api_rest"
	"github.com/deniskaponchik/GoSoft/internal/usecase/api_soap"
	"github.com/deniskaponchik/GoSoft/internal/usecase/authentication"
	"github.com/deniskaponchik/GoSoft/internal/usecase/repo"
	"github.com/deniskaponchik/GoSoft/pkg/postgres"
	"github.com/deniskaponchik/GoSoft/pkg/mysql1"
	"github.com/deniskaponchik/GoSoft/pkg/redis1"
	"os"
	"os/signal"
	"syscall"

	//"github.com/deniskaponchik/GoSoft/pkg/logger"
	//"github.com/deniskaponchik/GoSoft/pkg/rabbitmq/rmq_rpc/server"
	"log"
	"time"
)

//Run creates objects via constructors
func RunGisup(cfg *gisup.ConfigGisup) {

	//MySQL
	//repoGlpi, err := repo.NewRepoGlpi(cfg.GLPI.GlpiConnectStr, cfg.GLPI.DB)
	repoGlpi, err := mysql1.NewSqlMy(cfg.DbGlpi.GlpiConnectStr, cfg.DbGlpi.GlpiDB)
	if err != nil {
		log.Fatalf("Подключение к GLPI завершилось ошибкой: %w", err)
	} else {
		log.Println("Проверка подключения к БД GLPI прошла успешно")
	}
	//repoITsup, err := repo.NewRepoGlpi(cfg.GLPI.GlpiConnectStr, cfg.GLPI.DB)
	repoGisupMySqlProd, err := mysql1.NewSqlMy(cfg.DbGisupMySql.ITsupConnectStr, cfg.DbGisupMySql.ITsupDBprod)
	if err != nil {
		log.Fatalf("Подключение к GLPI завершилось ошибкой: %w", err)
	} else {
		log.Println("Проверка подключения к БД IT support прошла успешно")
	}
	repoGisupMySqlTest, err := mysql1.NewSqlMy(cfg.DbGisupMySql.ITsupConnectStr, cfg.DbGisupMySql.ITsupDBtest)
	if err != nil {
		log.Fatalf("Подключение к GLPI завершилось ошибкой: %w", err)
	} else {
		log.Println("Проверка подключения к БД IT support прошла успешно")
	}

	//Postgres
	repoPG, err := postgres.NewRepoPG(cfg.DbGisupPg.PgConnectStr)
	if err != nil {
		log.Printf("Подключение к Postgres завершилось ошибкой: %w", err)
	} else {
		log.Println("Проверка подключения к БД Postgres прошла успешно")
	}
	defer repoPG.Close()

	//Reddis
	repoRedis, err := redis1.NewRedis(cfg.DbRedis.RedisConnectString, cfg.DbRedis.RedisDB)
	if err != nil {
		log.Printf("Подключение к Redis завершилось ошибкой: %w", err)
	} else {
		log.Println("Проверка подключения к Redis прошла успешно")
	}

	repo, err := repo.NewRepo(
		*repoGlpi,
		*repoGisupMySqlProd,
		*repoGisupMySqlTest,
		*repoPG,
		*repoRedis,
	)

	/*SOAP
	gisupSoap, err := api_soap.NewSoap(cfg.SoapUrl, cfg.BpmUrl) // cfg.SoapTest, cfg.BpmTest
	if err != nil {
		log.Fatalf("Подключение к SOAP завершилось ошибкой: %w", err)
	} else {
		log.Println("Проверка подключения к SOAP прошла успешно")
	}*/
	soapTest := api_soap.NewSoap(cfg.SoapTest, cfg.BpmTest) /*
		if err != nil {
			log.Fatalf("Подключение к тестовому SOAP завершилось ошибкой: %w", err)
		} else {
			log.Println("Проверка подключения к тестовому SOAP прошла успешно")
		}*/
	soapProd := api_soap.NewSoap(cfg.SoapProd, cfg.BpmProd) /*
		if err != nil {
			log.Fatalf("Подключение к продуктовому SOAP завершилось ошибкой: %w", err)
		} else {
			log.Println("Проверка подключения к продуктовому SOAP прошла успешно")
		}*/
	log.Println("")

	//LDAP
	ldap := authentication.NewLdap(cfg.LdapDN, cfg.LdapDomain, "", "", cfg.LdapRoleDn, cfg.LdapServer)/*
		if err != nil {
			log.Println("Подключение к LDAP завершилось ошибкой: %w", err)
		} else {
			log.Println("Проверка подключения к SOAP прошла успешно")
		}*/
	log.Println("")

	//JWT
	jwt := authorization.NewAuthJwt(cfg.Token.JwtKey, cfg.Token.TTL)/*
		if err != nil {
			log.Println("Создание JWT-токена завершилось ошибкой: %w", err)
		} else {
			log.Println("Создание JWT-токена прошло успешно")
		}*/
	log.Println("")

	//C3PO
	c3po := api_rest.NewC3po(cfg.C3po.C3poUrl) /*
		err = c3po.GetUserLogin()
		if err != nil {
			log.Println("Подключение к C3PO завершилось ошибкой: %w", err)
		} else {
			log.Println("Проверка подключения к C3PO прошла успешно")
		}*/
	log.Println("")

	//RMQ client
	rmq := amqp_rmq.NewRmq(cfg.RmqConnectStr, cfg.RmqServerExchange)/*
		err = rmq.Publish("Start", "")
		if err != nil {
			log.Println("Подключение к RMQ завершилось ошибкой: %w", err)
		} else {
			log.Println("Проверка подключения к RMQ прошла успешно")
		}*/
	log.Println("")


	//
	gisupUseCase := usecase.NewGisupUC(
		//Обновление офисов из БД в мапу туда-обратно
		repo,
		c3po,
		//обработка админки веба
		jwt,
		ldap,
		//отправка в очередь общих сообщений всего приложения
		rmq,
		//очистка старых логов
		//
		soapProd,
		soapTest,

		)
	go gisupUseCase.InfinityProcessingGisup

	//
	wifiUseCase := usecase.NewWiFiuc(
		//Инициализурет мапу из БД
		gisupRepo,
		//создаёт заявки
		cfg.App.TimeZone,
		cfg.HttpURL,
		)
	go wifiUseCase.InfinityProcessingWiFi


		unifiUseCase := usecase.NewUnifiUC(
			//unifiRepo, //вставляем объект, который удовлетворяет интерфейсу UnifiRepo
			api_soap.NewSoap(cfg.SoapUrl, cfg.BpmUrl), // cfg.SoapTest, cfg.BpmTest
			amqp_rmq.NewRmqUnifi(cfg.RMQ.RmqConnectStr, cfg.RMQ.ServerExchange),

			api_web.NewUi(cfg.Ubiquiti.UiUsername, cfg.Ubiquiti.UiPassword, cfg.Ubiquiti.UiContrlRostov, 1),
			api_web.NewUi(cfg.Ubiquiti.UiUsername, cfg.Ubiquiti.UiPassword, cfg.Ubiquiti.UiContrlNovosib, 2),
			api_rest.NewUnifiC3po(cfg.C3po.C3poUrl)
			//authentication.NewLdap(cfg.LdapDN, cfg.LdapDomain, cfg.LdapLogin, cfg.LdapPassword, cfg.LdapRoleDn, cfg.LdapServer),
			//authorization.NewAuthJwt(cfg.Token.JwtKey, cfg.Token.TTL),

			cfg.App.EveryCodeMap,


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
		//gisupRepo,
		//gisupSoap,
	)

	vcsUseCase := usecase.NewVcsUc(
		//Инициализурет мапу из БД
		gisupRepo,
		//создаёт заявки
		cfg.App.TimeZone,
		cfg.HttpURL,
	)
	go vcsUseCase.InfinityTicketsCreating

	polyUseCase := usecase.NewPoly(
		gisupRepo, //polyRepo,
		//webapi.New(cfg.PolyUsername, cfg.PolyPassword),
		//api_web.NewPolyWebApi(cfg.PolyUsername, cfg.PolyPassword),
		api_rest.NewPolyWebApi(cfg.PolyUsername, cfg.PolyPassword),
		netdial.New(),
		api_soap.New(cfg.SoapUrl, cfg.BpmUrl), // cfg.SoapTest, cfg.BpmTest
		cfg.InnerVars.EveryCodeMap,
		cfg.InnerVars.RestartHour,
		cfg.App.TimeZone,
	)

	//GRPC
	myGrpc := myGRPC.New(
		//unifiUseCase,
		gisupUseCase,
		cfg.GrpcPort,
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
