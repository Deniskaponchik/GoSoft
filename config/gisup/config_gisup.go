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
	//100-заявки создаются минута в минуту без задержек по ночам
	timezone := flag.Int("timezone", 100, "Time hour from Moscow")

	//Базы данных будут называться одинаково, на каком бы сервере не находились
	//dbProd := flag.String("dbProd", "it_support_db_4", "database for gisup tables")
	//dbTest := flag.String("dbTest", "gisup_db_4", "database for gisup tables")
	dbGlpi := flag.String("db_glpi", "glpi_db", "database for glpi tables")
	dbGisup := flag.String("db_gisup", "it_support_db_4", "database for gisup mysql tables")
	dbPG := flag.String("db_pg", "gisup_db_4", "database for gisup postgres tables")
	dbRedis := flag.String("db_redis", "gisup_db_4", "database for gisup redis tables")

	//Аргументы внешних коммуникационных сервисов
	httpUrl := flag.String("http_url", "10.57.179.121:8081", "url of http-server")
	grpcPort := flag.Int("grpc_port", 8082, "port of grpc-server")
	rmqServExcahnge := flag.String("RmqServerExchange", "unifi.direct", "Exchange of rabbitMQ-server")
	tokenTTL := flag.Int("tokenTTL", 60, "minutes of live time token")

	//Unifi
	unifi_switch := flag.Int("unifi_switch", 1, "0 - Disable, 1 - Enable")
	unifi_mode := flag.String("unifi_mode", "TEST", "PROD, TEST, WEB")
	unifi_logLevel := flag.String("unifi_loglevel", "DEBUG", "DEBUG, INFO, WARN, ERROR")
	//Задаёт число месяца, когда заявки по аномалиям уже были созданы. Нужно для того, чтобы при старте кода не создавались
	//тикеты, а код ждал наступления следующего дня. Если необходимо, чтобы при старте кода шла проверка на аномалии, то
	//указать номер, отличный от сегодняшнего числа
	unifi_daily := flag.Int("unifi_daily", time.Now().Day(), "To do daily anomaly creating tickets")
	//TODO: вспомнить, за что это отвечает
	unifi_h1 := flag.Int("unifi_h1", time.Now().Hour(), "To do hourly anomaly downloading to DB")
	unifi_h2 := flag.Int("unifi_h2", time.Now().Hour(), "To do hourly anomaly downloading to DB")

	//Eltex
	eltex_switch := flag.Int("eltex_switch", 1, "0 - Disable, 1 - Enable")
	eltex_mode := flag.String("eltex_mode", "TEST", "PROD, TEST, WEB")
	eltex_logLevel := flag.String("eltex_loglevel", "DEBUG", "DEBUG, INFO, WARN, ERROR")

	//Polycom
	poly_switch := flag.Int("poly_switch", 1, "0 - Disable, 1 - Enable")
	poly_mode := flag.String("poly_mode", "TEST", "PROD, TEST, WEB")
	poly_logLevel := flag.String("poly_loglevel", "DEBUG", "DEBUG, INFO, WARN, ERROR")
	//чтобы отключить ежедневную перезагрузку, указать 25 и выше
	poly_restart_time := flag.Int("poly_restart_time", 7, "hour when codecs restart")

	//Lenovo
	lenovo_switch := flag.Int("lenovo_switch", 1, "0 - Disable, 1 - Enable")
	lenovo_mode := flag.String("lenovo_mode", "TEST", "PROD, TEST, WEB")
	lenovo_logLevel := flag.String("lenovo_loglevel", "DEBUG", "DEBUG, INFO, WARN, ERROR")

	flag.Parse()

	cfg.App.TimeZone = *timezone
	//cfg.Log.LevelCmd = *logLevel

	cfg.Token.TTL = *tokenTTL
	cfg.HttpURL = *httpUrl
	cfg.HttpPort = strings.Split(*httpUrl, ":")[1]
	cfg.GrpcPort = *grpcPort
	cfg.RmqServerExchange = *rmqServExcahnge

	cfg.Polycom.RestartHour = *poly_restart_time
	cfg.Ubiquiti.Daily = *unifi_daily
	cfg.Ubiquiti.H1 = *unifi_h1
	if cfg.Ubiquiti.H1 == 0 {
		cfg.Ubiquiti.H2 = 0
	} else {
		cfg.Ubiquiti.H2 = *unifi_h2
	}

	cfg.UiSwitch = *unifi_switch
	cfg.EltexSwitch = *eltex_switch
	cfg.PolySwitch = *poly_switch
	cfg.LenovoSwitch = *lenovo_switch

	cfg.UiMode = *unifi_mode
	cfg.EltexMode = *eltex_mode
	cfg.PolyMode = *poly_mode
	cfg.LenovoMode = *lenovo_mode

	cfg.UiLogLevel = *unifi_logLevel
	cfg.EltexLogLevel = *eltex_logLevel
	cfg.PolyLogLevel = *poly_logLevel
	cfg.LenovoLogLevel = *lenovo_logLevel

	/*
		cfg.GLPI.DB = *db
		cfg.DbGisupMySql.GisupDBprod = *dbProd
		cfg.DbGisupMySql.GisupDBprod = *dbTest
		cfg.DbGisupPg.PgDb = *dbTest  //тест, который в будущем станет продом
		cfg.DbRedis.RedisDB = *dbTest //тест, который в будущем станет продом
	*/
	cfg.DbGlpi.GlpiDB = *dbGlpi
	cfg.DbGisupMySql.GisupDB = *dbGisup
	cfg.DbGisupPg.PgDb = *dbPG
	cfg.DbRedis.RedisDB = *dbRedis

	//TODO: Нужно ли на данном этапе делать такое разделение?
	if *mode == "PROD" {
		//cfg.BpmUrl = cfg.BpmProd
		//cfg.SoapUrl = cfg.SoapProd

		wifiEveryCodeMap := map[int]int{
			//[минута]номер контроллера
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
		wifiEveryCodeMap := map[int]int{
			//[минута]номер контроллера
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
		cfg.Ubiquiti.UiEveryCodeMap = wifiEveryCodeMap
		cfg.Eltex.EltexEveryCodeMap = wifiEveryCodeMap
	}
	if *mode == "TEST" {
		//cfg.BpmUrl = cfg.BpmTest
		//cfg.SoapUrl = cfg.SoapTest

		//cfg.App.EveryCodeMap = map[int]int{ //[минута]номер контроллера
		wifiEveryCodeMap := map[int]int{ //[минута]номер контроллера
			5:  1,
			15: 2,
			25: 1,
			35: 2,
			45: 1,
			55: 2,
		}
		cfg.Ubiquiti.UiEveryCodeMap = wifiEveryCodeMap
		cfg.Eltex.EltexEveryCodeMap = wifiEveryCodeMap

		vcsEveryCodeMap := map[int]bool{
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
		cfg.Polycom.PolyEveryCodeMap = vcsEveryCodeMap
		cfg.Lenovo.LenovoEveryCodeMap = vcsEveryCodeMap

	}
	if *mode == "WEB" {
		//cfg.BpmUrl = cfg.BpmProd
		//cfg.SoapUrl = cfg.SoapProd

		//cfg.App.EveryCodeMap = make(map[int]int)
		cfg.Ubiquiti.UiEveryCodeMap = make(map[int]int)
		cfg.Eltex.EltexEveryCodeMap = make(map[int]int)

		cfg.Polycom.PolyEveryCodeMap = make(map[int]bool)
		cfg.Lenovo.LenovoEveryCodeMap = make(map[int]bool)

		cfg.Ubiquiti.UiSwitch = 0
		cfg.Eltex.EltexSwitch = 0
		cfg.Polycom.PolySwitch = 0
		cfg.Lenovo.LenovoSwitch = 0
	}

	log.Println("")

	log.Println("Mode	     : ", *mode) //cfg.InnerVars.Mode)
	//log.Println("Log level     : ", cfg.Log.LevelCmd)
	log.Println("Timezone    : ", cfg.App.TimeZone)
	log.Println("HTTP URL    : ", cfg.HttpURL)
	log.Println("C3PO URL    : ", cfg.C3poUrl)
	log.Println("")

	//log.Println("Every Code Map: ", cfg.App.EveryCodeMap)
	log.Println("Unifi Switch : ", cfg.Ubiquiti.UiSwitch)
	log.Println("Unifi Mode   : ", cfg.UiMode)
	log.Println("Unifi Log    : ", cfg.Ubiquiti.UiLogLevel)
	log.Println("Unifi Map    : ", cfg.Ubiquiti.UiEveryCodeMap)
	log.Println("")

	log.Println("Eltex Switch : ", cfg.Eltex.EltexSwitch)
	log.Println("Eltex Mode   : ", cfg.EltexMode)
	log.Println("Eltex Log    : ", cfg.Eltex.EltexLogLevel)
	log.Println("Eltex Map    : ", cfg.Eltex.EltexEveryCodeMap)
	log.Println("")

	log.Println("Poly Switch  : ", cfg.Polycom.PolySwitch)
	log.Println("Poly Mode    : ", cfg.PolyMode)
	log.Println("Poly Restart : ", cfg.Polycom.RestartHour)
	log.Println("")

	log.Println("Lenovo Switch: ", cfg.Lenovo.LenovoSwitch)
	log.Println("Lenovo Mode  : ", cfg.LenovoMode)
	log.Println("")

	return cfg, nil
}
