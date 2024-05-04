package gisup

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"strings"
	"time"
)

func NewConfigGisup() (*ConfigGisup, error) {
	cfg := &ConfigGisup{}

	//Подгрузка переменных с yaml файла.
	//err := cleanenv.ReadConfig("./config/config.yml", cfg) // в оригинале
	//err := cleanenv.ReadConfig("./config.yml", cfg)  //для тестирования
	//err := cleanenv.ReadConfig("../../../config.yml", cfg) // Unifi/cmd/poly/bin/Poly_v1.0
	//if err != nil {		return nil, log.Errorf("read config error: %w", err)	}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	//command line arguments
	//https://stackoverflow.com/questions/2707434/how-to-access-command-line-arguments-passed-to-a-go-program

	//Аргументы внешних коммуникационных сервисов
	httpUrl := flag.String("http", "10.57.179.121:8081", "url of http-server")
	grpcPort := flag.Int("grpc", 8082, "port of grpc-server")
	rmqServExcahnge := flag.String("RmqServExch", "unifi.direct", "Exchange of rabbitMQ-server")
	tokenTTL := flag.Int("tokenTTL", 60, "minutes of live time token")

	//Общие аргументы
	//share_mode влияет на то, какие будут использоваться сервера bpm и SOAP
	share_mode := flag.String("share_mode", "PROD", "mode of app work: PROD, TEST, WEB")
	//TODO: вспомнить, что тут можно указать и на что влияет
	share_timezone := flag.Int("share_timezone", 100, "Time hour from Moscow")
	//100-заявки создаются минута в минуту без задержек по ночам
	db := flag.String("db", "it_support_db_3", "database for unifi tables")

	//Unifi
	unifi_switch := flag.Int("unifi_switch", 1, "0 - Disable, 1 - Enable")
	unifi_daily := flag.Int("unifi_daily", time.Now().Day(), "To do daily anomaly creating tickets")
	//TODO: вспомнить, за что это отвечает
	unifi_h1 := flag.Int("unifi_h1", time.Now().Hour(), "To do hourly anomaly downloading to DB")
	unifi_h2 := flag.Int("unifi_h2", time.Now().Hour(), "To do hourly anomaly downloading to DB")
	//DEBUG, INFO, WARN, ERROR
	unifi_logLevel := flag.String("unifi_loglevel", "DEBUG", "level of log")

	//Eltex
	eltex_switch := flag.Int("eltex_switch", 1, "0 - Disable, 1 - Enable")
	//DEBUG, INFO, WARN, ERROR
	eltex_logLevel := flag.String("eltex_loglevel", "DEBUG", "level of log")

	//Polycom
	poly_switch := flag.Int("poly_switch", 1, "0 - Disable, 1 - Enable")
	//чтобы отключить ежедневную перезагрузку, указать 25 и выше
	poly_restart_time := flag.Int("poly_restart_time", 7, "hour when codecs restart")
	//DEBUG, INFO, WARN, ERROR
	poly_logLevel := flag.String("poly_loglevel", "DEBUG", "level of log")

	//AudioCodes
	audiocodes_switch := flag.Int("audiocodes_switch", 1, "0 - Disable, 1 - Enable")
	//DEBUG, INFO, WARN, ERROR
	audiocodes_logLevel := flag.String("audiocodes_loglevel", "DEBUG", "level of log")

	flag.Parse()

	if *share_mode == "TEST" {
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
	} else if *share_mode == "WEB" {
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

	cfg.GLPI.DB = *db
	cfg.App.TimeZone = *share_timezone
	//cfg.Log.LevelCmd = *logLevel

	cfg.Token.TTL = *tokenTTL
	cfg.HTTP.URL = *httpUrl
	cfg.HTTP.Port = strings.Split(*httpUrl, ":")[1]
	cfg.GRPC.Port = *grpcPort
	cfg.RMQ.ServerExchange = *rmqServExcahnge

	cfg.Ubiquiti.Daily = *unifi_daily
	cfg.Ubiquiti.H1 = *unifi_h1
	if cfg.Ubiquiti.H1 == 0 {
		cfg.Ubiquiti.H2 = 0
	} else {
		cfg.Ubiquiti.H2 = *unifi_h2
	}

	log.Println("Mode: ", *share_mode) //cfg.InnerVars.Mode)
	//log.Println("Log level: ", cfg.Log.LevelCmd)
	log.Println("Timezone: ", cfg.App.TimeZone)
	log.Println("HTTP URL: ", cfg.HTTP.URL)
	log.Println("C3PO URL: ", cfg.C3poUrl)
	log.Println("")

	log.Println("Every Code Map: ", cfg.App.EveryCodeMap)
	log.Println("")

	log.Println("")

	log.Println("")

	return cfg, nil
}
