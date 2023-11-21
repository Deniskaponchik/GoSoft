package app

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/config/ui"
	v1 "github.com/deniskaponchik/GoSoft/Unifi/internal/controller/http/v1"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/repo"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/soap"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/ubiq"
	"github.com/deniskaponchik/GoSoft/Unifi/pkg/httpserver"
	"github.com/deniskaponchik/GoSoft/Unifi/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Run creates objects via constructors.
func RunUnifi(cfg *ui.ConfigUi) {
	fmt.Println("")
	l := logger.New(cfg.Log.Level)

	/* Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()
	*/
	unifiRepo, err := repo.NewUnifiRepo(cfg.GLPI.GlpiITsupport, cfg.GLPI.GlpiConnectStrGLPI, cfg.UiContrlint)
	if err != nil {
		//если БД недоступна - останавливаем тут же
		l.Fatal(fmt.Errorf("app - Run - glpi.New: %w", err))
	} else {
		fmt.Println("Проверка подключения к БД прошла успешно")
	}

	//unpoller := unpoller.NewUnpoller(cfg.UiUsername, cfg.UiPassword, cfg.UiContrlstr)
	/*
		uc := unifi.Config{
			//c := *unifi.Config{  //ORIGINAL
			User: cfg.UiUsername,  //wifiConf.UnifiUsername,
			Pass: cfg.UiPassword,  //wifiConf.UnifiPassword,
			URL:  cfg.UiContrlstr, //urlController,
			// Log with log.Printf or make your own interface that accepts (msg, test_SOAP)
			ErrorLog: log.Printf,
			DebugLog: log.Printf,
		}*/

	unifiUseCase := usecase.NewUnifiUC(
		//repo.New(cfg.GLPI.GlpiConnectStrITsupport),
		unifiRepo,
		soap.NewSoap(cfg.SoapUrl, cfg.BpmUrl), // cfg.SoapTest, cfg.BpmTest
		//ubiq.NewUbiq(unpoller),                //cfg.Ubiquiti.UiUsername, cfg.Ubiquiti.UiPassword),
		ubiq.NewUi(cfg.Ubiquiti.UiUsername, cfg.Ubiquiti.UiPassword, cfg.Ubiquiti.UiContrlstr),
		cfg.App.EveryCodeMap,
		cfg.App.TimeZone,
	)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, unifiUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	err = unifiUseCase.InfinityProcessingUnifi() //cfg.BpmUrl, cfg.SoapUrl)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - InfinityUnifiProcessing: %w", err))
	}

}
