package ui

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

func NewConfigUnifi() (*ConfigUi, error) {
	cfg := &ConfigUi{}

	//Подгрузка переменных с yaml файла. Отключаю из-за геморроя с указанием пути до него
	//err := cleanenv.ReadConfig("./config/config.yml", cfg) // в оригинале
	//err := cleanenv.ReadConfig("./config.yml", cfg)  //для тестирования
	//err := cleanenv.ReadConfig("../../../config.yml", cfg) // Unifi/cmd/poly/bin/Poly_v1.0
	//if err != nil {		return nil, fmt.Errorf("read config error: %w", err)	}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	everyCodeSlice := [21]map[int]bool{}
	//everyone 6 minutes, between = 3 minutes, Start at Rostov = 3, Novosib = 6
	everyCodeSlice[1] = map[int]bool{
		3:  true,
		9:  true,
		15: true,
		21: true,
		27: true,
		33: true,
		39: true,
		45: true,
		51: true,
		57: true,
	}
	everyCodeSlice[11] = map[int]bool{
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

	//everyone 12 minutes, between = 6 minutes, Start at Rostov = 9, Novosib = 3
	everyCodeSlice[2] = map[int]bool{
		9:  true,
		21: true,
		33: true,
		45: true,
		57: true,
	}
	everyCodeSlice[12] = map[int]bool{
		3:  true,
		15: true,
		27: true,
		39: true,
		51: true,
	}

	//everyone 20 minutes, between = 10 minutes, Start at Rostov = 5, Novosib = 15
	everyCodeSlice[3] = map[int]bool{
		5:  true,
		25: true,
		45: true,
	}
	everyCodeSlice[13] = map[int]bool{
		15: true,
		35: true,
		55: true,
	}

	//https://stackoverflow.com/questions/2707434/how-to-access-command-line-arguments-passed-to-a-go-program
	//mode := "TEST"
	mode := flag.String("mode", "PROD", "mode of app work: PROD, TEST")
	controller := flag.String("cntrl", "Rostov", "controller: Novosib, Rostov")
	flag.Parse()

	//cfg.InnerVars.Mode = *mode
	if *mode == "TEST" {
		cfg.BpmUrl = cfg.BpmTest
		cfg.SoapUrl = cfg.SoapTest
		//cfg.GlpiITsupport = cfg.GlpiITsupportTest
		cfg.GlpiITsupport = "root:t2root@tcp(10.77.252.153:3306)/it_support_test_db"
	} else {
		// "PROD"
		cfg.BpmUrl = cfg.BpmProd
		cfg.SoapUrl = cfg.SoapProd
		cfg.GlpiITsupport = cfg.GlpiITsupportProd
	}

	//controller = *controller //
	if *controller == "Rostov" {
		cfg.Ubiquiti.UiContrlstr = cfg.Ubiquiti.UiContrlRostov
		cfg.Ubiquiti.UiContrlint = 1
		cfg.App.EveryCodeMap = everyCodeSlice[2] //каждые 12 минут
	} else {
		// "Novosib"
		cfg.Ubiquiti.UiContrlstr = cfg.Ubiquiti.UiContrlNovosib
		cfg.Ubiquiti.UiContrlint = 2
		cfg.App.EveryCodeMap = everyCodeSlice[12] //каждые 12 минут
	}

	fmt.Println("Mode: ", *mode) //cfg.InnerVars.Mode)
	fmt.Println("Controller: ", cfg.Ubiquiti.UiContrlstr)
	fmt.Println("Every Code Map: ", cfg.App.EveryCodeMap)
	//time.Sleep(1000 * time.Second)

	return cfg, nil
}

type (
	ConfigUi struct {
		//Polycom
		Ubiquiti
		Bpm
		Soap
		GLPI

		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		PG   `yaml:"postgres"`
		RMQ  `yaml:"rabbitmq"`
	}

	App struct {
		Name         string `yaml:"name"`
		Version      string `yaml:"version"`
		EveryCodeMap map[int]bool
	}
	//env-required:"true" - ОБЯЗАТЕЛЬНО должен получить перменную либо из окружения, либо из yaml. Между true и false разницы не заметил

	/*
		Polycom struct {
			PolyUsername string `env-required:"true" yaml:"poly_usernamename"    env:"POLY_USERNAME"`
			PolyPassword string `env-required:"true" yaml:"poly_password"        env:"POLY_PASSWORD"`
		}*/
	Ubiquiti struct {
		UiUsername      string `env-required:"true" yaml:"unifi_usernamename"   env:"UNIFI_USERNAME"`
		UiPassword      string `env-required:"true" yaml:"unifi_password"       env:"UNIFI_PASSWORD"`
		UiContrlRostov  string `env-required:"true" yaml:"contrl_rostov"   env:"UNIFI_CONTROLLER_ROSTOV"`
		UiContrlNovosib string `env-required:"true" yaml:"contrl_novosib"  env:"UNIFI_CONTROLLER_NOVOSIB"`
		UiContrlstr     string
		UiContrlint     int
	}
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
		GlpiITsupportTest  string //`env-required:"true"   env:"GLPI_ITSUP_TEST"`
		GlpiITsupport      string //`env-required:"false"`
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
