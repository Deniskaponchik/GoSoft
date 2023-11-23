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
	"os"
	"os/signal"
	"syscall"
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

	unifiUseCase := usecase.NewUnifiUC(
		//repo.New(cfg.GLPI.GlpiConnectStrITsupport),
		unifiRepo,
		soap.NewSoap(cfg.SoapUrl, cfg.BpmUrl), // cfg.SoapTest, cfg.BpmTest
		//ubiq.NewUbiq(unpoller),                //cfg.Ubiquiti.UiUsername, cfg.Ubiquiti.UiPassword),
		ubiq.NewUi(cfg.Ubiquiti.UiUsername, cfg.Ubiquiti.UiPassword, cfg.Ubiquiti.UiContrlstr),
		cfg.App.EveryCodeMap,
		cfg.App.TimeZone,
	)

	go unifiUseCase.InfinityProcessingUnifi()
	fmt.Println("InfinityProcessingUnifi отправился в горутину")
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

	// HTTP Server
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
	}
}
