package poly

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

func NewConfigPoly() (*ConfigPoly, error) {
	cfg := &ConfigPoly{}

	//Подгрузка переменных с yaml файла. Отключаю из-за геморроя с указанием пути до него
	//err := cleanenv.ReadConfig("./config/config.yml", cfg) // в оригинале
	//err := cleanenv.ReadConfig("./config.yml", cfg)  //для тестирования
	//err := cleanenv.ReadConfig("../../../config.yml", cfg) // Unifi/cmd/poly/bin/Poly_v1.0
	//if err != nil {		return nil, fmt.Errorf("read config error: %w", err)	}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	everyCodeSlice := [4]map[int]bool{}
	//every 20 minutes, start at 5
	everyCodeSlice[0] = map[int]bool{
		5:  true,
		25: true,
		45: true,
	}
	//every 6 minutes, run at 6
	everyCodeSlice[1] = map[int]bool{
		6:  true,
		12: true,
		18: true,
		24: true,
		30: true,
		36: true,
		42: true,
		48: true,
		54: true,
		59: true,
	}
	//every 6 minutes, run at 3
	everyCodeSlice[2] = map[int]bool{
		3:  true,
		9:  true,
		15: true,
		21: true,
		33: true,
		39: true,
		45: true,
		51: true,
		57: true,
	}
	//every 3 minutes, run at 3. Without 00:00
	everyCodeSlice[3] = map[int]bool{
		3:  true,
		6:  true,
		9:  true,
		12: true,
		15: true,
		18: true,
		21: true,
		24: true,
		27: true,
		30: true,
		33: true,
		36: true,
		39: true,
		42: true,
		45: true,
		48: true,
		51: true,
		54: true,
		57: true,
	}

	//https://stackoverflow.com/questions/2707434/how-to-access-command-line-arguments-passed-to-a-go-program
	//mode := "TEST"
	mode := flag.String("mode", "PROD", "mode of app work: PROD, TEST")
	restart := flag.Int("restart", 7, "hour when codecs restart") //чтобы отключить ежедневную перезагрузку, указать 25 и выше
	timezone := flag.Int("time", 100, "Time hour from Moscow")    //100-заявки создаются минута в минуту без задержек по ночам
	flag.Parse()

	cfg.InnerVars.Mode = *mode
	if cfg.InnerVars.Mode == "TEST" {
		cfg.BpmUrl = cfg.BpmTest
		cfg.SoapUrl = cfg.SoapTest
		cfg.GlpiITsupport = cfg.GlpiITsupportProd
		cfg.InnerVars.EveryCodeMap = everyCodeSlice[3] //каждые 3 минут
	} else {
		// "PROD"
		cfg.BpmUrl = cfg.BpmProd
		cfg.SoapUrl = cfg.SoapProd
		cfg.GlpiITsupport = cfg.GlpiITsupportTest
		cfg.InnerVars.EveryCodeMap = everyCodeSlice[0] //каждые 20 минут, старт в 5 минут
	}
	cfg.InnerVars.RestartHour = *restart //в 7 часов по времени сервера, где запущен скрипт
	cfg.App.TimeZone = *timezone

	fmt.Println("Mode: ", cfg.InnerVars.Mode)
	fmt.Println("Restart hour: ", cfg.InnerVars.RestartHour)
	fmt.Println("Every Code Map: ", cfg.InnerVars.EveryCodeMap)
	fmt.Println("Timezone: ", cfg.App.TimeZone)
	//time.Sleep(1000 * time.Second)

	return cfg, nil
}

type (
	ConfigPoly struct {
		InnerVars

		Polycom
		//Ubiquiti
		Bpm
		Soap
		GLPI

		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		PG   `yaml:"postgres"`
		RMQ  `yaml:"rabbitmq"`
	}
	InnerVars struct {
		Mode         string
		EveryCodeMap map[int]bool
		RestartHour  int
		//EveryCodeSlice [4]map[int]bool
	}
	//env-required:"true" - ОБЯЗАТЕЛЬНО должен получить перменную либо из окружения, либо из yaml. Между true и false разницы не заметил

	Polycom struct {
		PolyUsername string `env-required:"true" yaml:"poly_usernamename"    env:"POLY_USERNAME"`
		PolyPassword string `env-required:"true" yaml:"poly_password"        env:"POLY_PASSWORD"`
	}
	/*
		Ubiquiti struct {
			UiUsername      string `env-required:"true" yaml:"unifi_usernamename"   env:"UNIFI_USERNAME"`
			UiPassword      string `env-required:"true" yaml:"unifi_password"       env:"UNIFI_PASSWORD"`
			UiContrlRostov  string `env-required:"true" yaml:"contrl_rostov"   env:"UNIFI_CONTROLLER_ROSTOV"`
			UiContrlNovosib string `env-required:"true" yaml:"contrl_novosib"  env:"UNIFI_CONTROLLER_NOVOSIB"`
		}*/
	Bpm struct {
		BpmUrl  string //`env-required:"false"`
		BpmProd string `env-required:"true" yaml:"bpm_prod"   env:"BPM_PROD"`
		BpmTest string `env-required:"true" yaml:"bpm_test"   env:"BPM_TEST"`
	}
	Soap struct {
		SoapUrl  string //`env-required:"false"`
		SoapProd string `env-required:"true" env:"SOAP_PROD"`
		SoapTest string `env-required:"true" env:"SOAP_TEST"`
	}
	GLPI struct {
		GlpiConnectStrGLPI string `env-required:"true"   env:"GLPI_CONNECT_STR_GLPI"`
		GlpiITsupportProd  string `env-required:"true"   env:"GLPI_CONNECT_STR_ITSUP"`
		GlpiITsupportTest  string `env-required:"true"   env:"GLPI_ITSUP_TEST"`
		GlpiITsupport      string //`env-required:"false"`
	}

	App struct {
		Name     string `yaml:"name"`
		Version  string `yaml:"version"`
		TimeZone int
	}
	Log struct {
		Level string `yaml:"log_level"   env:"LOG_LEVEL"`
	}
	HTTP struct {
		Port string `yaml:"port" env:"HTTP_PORT"`
	}
	PG struct {
		PoolMax int `yaml:"pool_max" env:"PG_POOL_MAX"`
		//URL     string `env-required:"true"                 env:"PG_URL"`
	}
	RMQ struct {
		ServerExchange string `yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
		ClientExchange string `yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
		//URL            string `env-required:"true"                            env:"RMQ_URL"`
	}
)
