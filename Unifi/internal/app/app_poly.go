package app

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/Unifi/config/poly"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/netdial"
	_ "github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/ping"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/repo"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/soap"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/usecase/webapi"
	"github.com/deniskaponchik/GoSoft/Unifi/pkg/logger"
)

// Run creates objects via constructors.
func PolyRun(cfg *poly.ConfigPoly) {
	fmt.Println("")
	l := logger.New(cfg.Log.Level)

	/* Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()
	*/
	polyRepo, err := repo.New(cfg.GLPI.GlpiITsupport)
	if err != nil {
		//если БД недоступна - останавливаем тут же
		l.Fatal(fmt.Errorf("app - Run - glpi.New: %w", err))
	} else {
		fmt.Println("Проверка подключения к БД прошла успешно")
	}

	polyUseCase := usecase.NewPoly(
		//repo.New(cfg.GLPI.GlpiConnectStrITsupport),
		polyRepo,
		webapi.New(cfg.PolyUsername, cfg.PolyPassword),
		netdial.New(),
		soap.New(cfg.SoapUrl, cfg.BpmUrl), // cfg.SoapTest, cfg.BpmTest
		cfg.InnerVars.EveryCodeMap,
		cfg.InnerVars.RestartHour,
	)

	err = polyUseCase.InfinityPolyProcessing() //cfg.BpmUrl, cfg.SoapUrl)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - InfinityPolyProcessing: %w", err))
	}

}
