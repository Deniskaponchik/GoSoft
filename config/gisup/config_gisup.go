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

	//Общие аргументы
	//share_mode влияет на то, какие будут использоваться сервера bpm и SOAP
	mode := flag.String("mode", "PROD", "mode of app work: PROD, TEST, WEB")
	//TODO: вспомнить, что тут можно указать и на что влияет
	timezone := flag.Int("timezone", 100, "Time hour from Moscow")
	//100-заявки создаются минута в минуту без задержек по ночам
	db := flag.String("db", "it_support_db_3", "database for unifi tables")

	//Аргументы внешних коммуникационных сервисов
	httpUrl := flag.String("http_url", "10.57.179.121:8081", "url of http-server")
	grpcPort := flag.Int("grpc_port", 8082, "port of grpc-server")
	rmqServExcahnge := flag.String("RmqServerExchange", "unifi.direct", "Exchange of rabbitMQ-server")
	tokenTTL := flag.Int("tokenTTL", 60, "minutes of live time token")

	//Unifi
	unifi_switch := flag.Int("unifi_switch", 1, "0 - Disable, 1 - Enable")
	//Задаёт число месяца, когда заявки по аномалиям уже были созданы. Нужно для того, чтобы при старте кода не создавались
	//тикеты, а код ждал наступления следующего дня. Если необходимо, чтобы при старте кода шла проверка на аномалии, то
	//указать номер, отличный от сегодняшнего числа
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

	cfg.GLPI.DB = *db
	cfg.App.TimeZone = *timezone
	//cfg.Log.LevelCmd = *logLevel

	cfg.Token.TTL = *tokenTTL
	cfg.HTTP.URL = *httpUrl
	cfg.HTTP.Port = strings.Split(*httpUrl, ":")[1]
	cfg.GRPC.Port = *grpcPort
	cfg.RMQ.ServerExchange = *rmqServExcahnge

	cfg.Polycom.RestartHour = *poly_restart_time
	cfg.Ubiquiti.Daily = *unifi_daily
	cfg.Ubiquiti.H1 = *unifi_h1
	if cfg.Ubiquiti.H1 == 0 {
		cfg.Ubiquiti.H2 = 0
	} else {
		cfg.Ubiquiti.H2 = *unifi_h2
	}

	cfg.Ubiquiti.UiSwitch = *unifi_switch
	cfg.Eltex.EltexSwitch = *eltex_switch
	cfg.Polycom.PolySwitch = *poly_switch
	cfg.AudioCodes.AudioSwitch = *audiocodes_switch

	cfg.Ubiquiti.UiLogLevel = *unifi_logLevel
	cfg.Eltex.EltexLogLevel = *eltex_logLevel
	cfg.Polycom.PolyLogLevel = *poly_logLevel
	cfg.AudioCodes.AudioLogLevel = *audiocodes_logLevel

	if *mode == "TEST" {
		cfg.BpmUrl = cfg.BpmTest
		cfg.SoapUrl = cfg.SoapTest

		//cfg.App.EveryCodeMap = map[int]int{ //[минута]номер контроллера
		shareEveryCodeMap := map[int]int{ //[минута]номер контроллера
			5:  1,
			15: 2,
			25: 1,
			35: 2,
			45: 1,
			55: 2,
		}
		cfg.Ubiquiti.UiEveryCodeMap = shareEveryCodeMap
		cfg.Eltex.EltexEveryCodeMap = shareEveryCodeMap

	} else if *mode == "WEB" {
		cfg.BpmUrl = cfg.BpmProd
		cfg.SoapUrl = cfg.SoapProd

		//cfg.App.EveryCodeMap = make(map[int]int)
		cfg.Ubiquiti.UiEveryCodeMap = make(map[int]int)
		cfg.Eltex.EltexEveryCodeMap = make(map[int]int)

		cfg.Ubiquiti.UiSwitch = 0
		cfg.Eltex.EltexSwitch = 0
		cfg.Polycom.PolySwitch = 0
		cfg.AudioCodes.AudioSwitch = 0

	} else if *mode == "PROD" {
		cfg.BpmUrl = cfg.BpmProd
		cfg.SoapUrl = cfg.SoapProd

		//cfg.App.EveryCodeMap = map[int]int{ //[минута]номер контроллера
		shareEveryCodeMap := map[int]int{ //[минута]номер контроллера
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
		cfg.Ubiquiti.UiEveryCodeMap = shareEveryCodeMap
		cfg.Eltex.EltexEveryCodeMap = shareEveryCodeMap

	}
	log.Println("")

	log.Println("Mode	     : ", *mode) //cfg.InnerVars.Mode)
	//log.Println("Log level     : ", cfg.Log.LevelCmd)
	log.Println("Timezone    : ", cfg.App.TimeZone)
	log.Println("HTTP URL    : ", cfg.HTTP.URL)
	log.Println("C3PO URL    : ", cfg.C3poUrl)
	log.Println("")

	//log.Println("Every Code Map: ", cfg.App.EveryCodeMap)
	log.Println("Unifi Switch: ", cfg.Ubiquiti.UiSwitch)
	log.Println("Unifi Log   : ", cfg.Ubiquiti.UiLogLevel)
	log.Println("Unifi Map   : ", cfg.Ubiquiti.UiEveryCodeMap)
	log.Println("")

	log.Println("Eltex Switch: ", cfg.Eltex.EltexSwitch)
	log.Println("Eltex Log   : ", cfg.Eltex.EltexLogLevel)
	log.Println("Eltex Map   : ", cfg.Eltex.EltexEveryCodeMap)
	log.Println("")

	log.Println("Poly Switch : ", cfg.Polycom.PolySwitch)
	log.Println("Poly Restart: ", cfg.Polycom.RestartHour)
	log.Println("")

	log.Println("Audio Switch: ", cfg.AudioCodes.AudioSwitch)
	log.Println("")

	return cfg, nil
}
