package app

import (
	"fmt"
	"github.com/deniskaponchik/GoSoft/config/poly"
	"github.com/deniskaponchik/GoSoft/internal/usecase"
	"github.com/deniskaponchik/GoSoft/internal/usecase/api_rest"
	"github.com/deniskaponchik/GoSoft/internal/usecase/api_soap"
	//"github.com/deniskaponchik/GoSoft/internal/usecase/api_web"
	"github.com/deniskaponchik/GoSoft/internal/usecase/netdial"
	"github.com/deniskaponchik/GoSoft/internal/usecase/repo"
	"github.com/deniskaponchik/GoSoft/pkg/logger"
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
	//polyRepo, err := repo.NewPolyRepo(cfg.GLPI.GlpiITsupport)
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

	err = polyUseCase.InfinityPolyProcessing() //cfg.BpmUrl, cfg.SoapUrl)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - InfinityPolyProcessing: %w", err))
	}

}
