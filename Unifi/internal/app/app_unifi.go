package app

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/config/ui"
	//"github.com/deniskaponchik/GoSoft/Unifi/internal/controller/http/fokusov"
	fokInterface "github.com/deniskaponchik/GoSoft/Unifi/internal/controller/http/fokInterface"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/repo"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/soap"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/ubiq"
	"log"
	//"github.com/deniskaponchik/GoSoft/Unifi/pkg/logger"
)

// Run creates objects via constructors.
func RunUnifi(cfg *ui.ConfigUi) {
	//fmt.Println("")

	//удалить префикс времени в логах
	//https://stackoverflow.com/questions/48629988/remove-timestamp-prefix-from-go-logger
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetFlags(0)
	//log.SetOutput()
	log.Println("")

	//Zerro Log  	//l := logger.New(cfg.Log.Level)   //l.Info("")

	/* Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()
	*/
	//unifiRepo, err := repo.NewUnifiRepo(cfg.GLPI.GlpiITsupport, cfg.GLPI.GlpiConnectStrGLPI) //, cfg.UiContrlint)
	unifiRepo, err := repo.NewUnifiRepo(cfg.GLPI.GlpiConnectStr, cfg.GLPI.DB)
	//repoRostov, err := repo.NewUnifiRepo(cfg.GLPI.GlpiITsupport, cfg.GLPI.GlpiConnectStrGLPI, cfg.)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - glpi.New: %w", err))
	} else {
		log.Println("Проверка подключения к БД прошла успешно")
	}

	unifiUseCase := usecase.NewUnifiUC(
		//repo.New(cfg.GLPI.GlpiConnectStrITsupport),
		unifiRepo,                             //вставляем объект, который удовлетворяет интерфейсу UnifiRepo
		soap.NewSoap(cfg.SoapUrl, cfg.BpmUrl), // cfg.SoapTest, cfg.BpmTest
		//ubiq.NewUbiq(unpoller),                //cfg.Ubiquiti.UiUsername, cfg.Ubiquiti.UiPassword),
		//ubiq.NewUi(cfg.Ubiquiti.UiUsername, cfg.Ubiquiti.UiPassword, cfg.Ubiquiti.UiContrlstr),
		ubiq.NewUi(cfg.Ubiquiti.UiUsername, cfg.Ubiquiti.UiPassword, cfg.Ubiquiti.UiContrlRostov, 1),
		ubiq.NewUi(cfg.Ubiquiti.UiUsername, cfg.Ubiquiti.UiPassword, cfg.Ubiquiti.UiContrlNovosib, 2),
		cfg.App.EveryCodeMap,
		cfg.App.TimeZone,
		cfg.HTTP.URL,
	)

	go unifiUseCase.InfinityProcessingUnifi()
	//l.Info("InfinityProcessingUnifi отправился в горутину")
	log.Println("InfinityProcessingUnifi отправился в горутину")
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
	//первоначальная версия:
	//err = unifiUseCase.InfinityProcessingUnifi() //cfg.BpmUrl, cfg.SoapUrl)
	//if err != nil {		l.Fatal(fmt.Errorf("app - Run - InfinityUnifiProcessing: %w", err))	}

	//FOKUSOV
	//router := *gin.Engine
	//httpFokusov := fokusov.New(
	httpFokusov := fokInterface.New(
		//gin.Engine,
		unifiUseCase,
		//usecase.Rest(),
		cfg.HTTP.Port,
	)
	httpFokusov.Start()

	/* EVRONE
	handler := gin.New()
	v1.NewRouter(handler, l, unifiUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}*/
}
