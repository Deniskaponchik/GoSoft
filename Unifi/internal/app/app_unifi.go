package app

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/config/ui"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/repo"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/soap"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/unifi"
	"github.com/deniskaponchik/GoSoft/Unifi/pkg/logger"
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
	unifiRepo, err := repo.NewUnifiRepo(cfg.GLPI.GlpiITsupportTest, cfg.GLPI.GlpiConnectStrGLPI, cfg.UiContrlint)
	if err != nil {
		//если БД недоступна - останавливаем тут же
		l.Fatal(fmt.Errorf("app - Run - glpi.New: %w", err))
	} else {
		fmt.Println("Проверка подключения к БД прошла успешно")
	}

	unifiUseCase := usecase.NewUnifi(
		//repo.New(cfg.GLPI.GlpiConnectStrITsupport),
		unifiRepo,
		soap.New(cfg.SoapUrl, cfg.BpmUrl), // cfg.SoapTest, cfg.BpmTest
		unifi.New(cfg.Ubiquiti.UiUsername, cfg.Ubiquiti.UiPassword),
		cfg.App.EveryCodeMap,
	)

	err = unifiUseCase.InfinityProcessingUnifi() //cfg.BpmUrl, cfg.SoapUrl)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - InfinityPolyProcessing: %w", err))
	}

}
