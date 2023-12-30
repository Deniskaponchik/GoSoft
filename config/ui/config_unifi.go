package ui

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"strings"
	"time"
)

func NewConfigUnifi() (*ConfigUi, error) {
	cfg := &ConfigUi{}

	//Подгрузка переменных с yaml файла.
	//err := cleanenv.ReadConfig("./config/config.yml", cfg) // в оригинале
	//err := cleanenv.ReadConfig("./config.yml", cfg)  //для тестирования
	//err := cleanenv.ReadConfig("../../../config.yml", cfg) // Unifi/cmd/poly/bin/Poly_v1.0
	//if err != nil {		return nil, log.Errorf("read config error: %w", err)	}

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

	//command line arguments
	//https://stackoverflow.com/questions/2707434/how-to-access-command-line-arguments-passed-to-a-go-program
	mode := flag.String("mode", "PROD", "mode of app work: PROD, TEST")
	db := flag.String("db", "it_support_db_3", "database for unifi tables")
	//controller := flag.String("cntrl", "Rostov", "controller: Novosib, Rostov")
	timezone := flag.Int("time", 100, "Time hour from Moscow") //100-заявки создаются минута в минуту без задержек по ночам
	httpUrl := flag.String("httpUrl", "wsir-it-03:8081", "url of web-server")
	daily := flag.Int("daily", time.Now().Day(), "To do daily anomaly creating tickets")
	h1 := flag.Int("h1", time.Now().Hour(), "To do hourly anomaly downloading to DB")
	h2 := flag.Int("h2", time.Now().Hour(), "To do hourly anomaly downloading to DB")
	logLevel := flag.String("log", "DEBUG", "level of log") //DEBUG, INFO, WARN, ERROR
	flag.Parse()

	/*
		cfg.App.EveryCodeMap = map[int]int{ //[минута]номер контроллера
			2:  2, // в начале часа различные выгрузки/загрузки в БД. нужно больше времени
			9:  1,
			15: 2,
			21: 1,
			27: 2,
			33: 1,
			39: 2,
			45: 1,
			51: 2,
			57: 1,
		}*/
	if *mode == "TEST" {
		cfg.BpmUrl = cfg.BpmTest
		cfg.SoapUrl = cfg.SoapTest
		cfg.App.EveryCodeMap = map[int]int{ //[минута]номер контроллера
			5:  1,
			15: 2,
			25: 1,
			35: 2,
			45: 1,
			55: 2,
		}
	} else if *mode == "WEB" {
		cfg.BpmUrl = cfg.BpmProd
		cfg.SoapUrl = cfg.SoapProd
		cfg.App.EveryCodeMap = make(map[int]int)
	} else {
		// "PROD"
		cfg.BpmUrl = cfg.BpmProd
		cfg.SoapUrl = cfg.SoapProd
		cfg.App.EveryCodeMap = map[int]int{ //[минута]номер контроллера
			2:  1, // в начале часа различные выгрузки/загрузки в БД. нужно больше времени
			9:  2,
			15: 1,
			21: 2,
			27: 1,
			33: 2,
			39: 1,
			45: 2,
			51: 1,
			57: 2,
		}
	}

	/*controller = *controller //
	if *controller == "Rostov" {
		cfg.Ubiquiti.UiContrlstr = cfg.Ubiquiti.UiContrlRostov
		cfg.Ubiquiti.UiContrlint = 1
		cfg.App.EveryCodeMap = everyCodeSlice[2] //каждые 12 минут
	} else {
		// "Novosib"
		cfg.Ubiquiti.UiContrlstr = cfg.Ubiquiti.UiContrlNovosib
		cfg.Ubiquiti.UiContrlint = 2
		cfg.App.EveryCodeMap = everyCodeSlice[12] //каждые 12 минут
	}*/
	cfg.GLPI.DB = *db
	cfg.App.TimeZone = *timezone
	cfg.HTTP.URL = *httpUrl
	cfg.HTTP.Port = strings.Split(*httpUrl, ":")[1]
	cfg.Ubiquiti.Daily = *daily
	cfg.Ubiquiti.H1 = *h1
	if cfg.Ubiquiti.H1 == 0 {
		cfg.Ubiquiti.H2 = 0
	} else {
		cfg.Ubiquiti.H2 = *h2
	}
	cfg.Log.LevelCmd = *logLevel

	log.Println("Mode: ", *mode) //cfg.InnerVars.Mode)
	log.Println("Log level: ", cfg.Log.LevelCmd)
	log.Println("Every Code Map: ", cfg.App.EveryCodeMap)
	log.Println("Timezone: ", cfg.App.TimeZone)
	log.Println("HTTP URL: ", cfg.HTTP.URL)
	log.Println("C3PO URL: ", cfg.C3poUrl)

	return cfg, nil
}

type (
	ConfigUi struct {
		//Polycom
		Ubiquiti
		Bpm
		Soap
		GLPI
		C3po
		Ldap

		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		//PG   `yaml:"postgres"`
		//RMQ  `yaml:"rabbitmq"`
	}
	App struct {
		Name         string `yaml:"name"`
		Version      string `yaml:"version"`
		EveryCodeMap map[int]int
		TimeZone     int
	}
	//env-required:"true" -ОБЯЗАТЕЛЬНО должен получить переменную либо из окружения, либо из yaml. Между true и false разницы не заметил

	/*Polycom struct {
		PolyUsername string `env-required:"true" yaml:"poly_usernamename"    env:"POLY_USERNAME"`
		PolyPassword string `env-required:"true" yaml:"poly_password"        env:"POLY_PASSWORD"`
	}*/
	Ubiquiti struct {
		UiUsername      string `env-required:"true" yaml:"unifi_usernamename"   env:"UNIFI_USERNAME"`
		UiPassword      string `env-required:"true" yaml:"unifi_password"       env:"UNIFI_PASSWORD"`
		UiContrlRostov  string `env-required:"true" yaml:"contrl_rostov"   env:"UNIFI_CONTROLLER_ROSTOV"`
		UiContrlNovosib string `env-required:"true" yaml:"contrl_novosib"  env:"UNIFI_CONTROLLER_NOVOSIB"`
		//UiContrlstr     string
		//UiContrlint     int //для совместного приложения двух контроллеров не должен приходить с конфигом
		Daily int
		H1    int
		H2    int
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
		GlpiConnectStr string `env-required:"true"   env:"GLPI_CONNECT_STR"` //строка подключения к серверу без указания БД
		//GlpiConnectStrGLPI string `env-required:"true"   env:"GLPI_CONNECT_STR_GLPI"`
		//GlpiITsupportProd  string `env-required:"true"   env:"GLPI_CONNECT_STR_ITSUP"`
		//GlpiITsupportTest  string `env-required:"true"   env:"GLPI_ITSUP_TEST"`
		//GlpiITsupport      string //`env-required:"false"`
		DB string //имя базы данных для unifi таблиц. задаю аргументами командной строки
	}
	C3po struct {
		C3poUrl string `env-required:"true"   env:"C3PO_URL"`
	}
	Ldap struct {
		LdapDN       string `env-required:"true"   env:"LDAP_DN"`
		LdapDomain   string `env-required:"true"   env:"LDAP_Domain"`
		LdapLogin    string `env-required:"true"   env:"LDAP_LOGIN"`
		LdapPassword string `env-required:"true"   env:"LDAP_PASSWORD"`
		LdapRoleDn   string `env-required:"true"   env:"LDAP_ROLE_DN"`
		LdapServer   string `env-required:"true"   env:"LDAP_SERVER"`
	}

	Log struct {
		LevelEnv string `yaml:"log_level"   env:"LOG_LEVEL"`
		LevelCmd string
	}
	HTTP struct {
		URL    string
		Port   string //`yaml:"port" env:"HTTP_PORT"`
		JwtKey string `env-required:"true"   env:"JWT_KEY"`
	}
	/*
		PG struct {
			PoolMax int `yaml:"pool_max" env:"PG_POOL_MAX"`
			//URL     string `env-required:"true"                 env:"PG_URL"`
		}
		RMQ struct {
			ServerExchange string `yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
			ClientExchange string `yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
			//URL            string `env-required:"true"                            env:"RMQ_URL"`
		}*/
)
