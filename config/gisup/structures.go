package gisup

type (
	ConfigGisup struct {
		Polycom
		Lenovo
		Ubiquiti
		Eltex
		Bpm
		Soap
		DbGlpi
		DbGisupMySql
		DbGisupPg //`yaml:"postgres"`
		DbRedis
		C3po
		Ldap

		//Log `yaml:"logger"`
		App `yaml:"app"`
		Token
		HTTP `yaml:"http"`
		GRPC
		RMQ `yaml:"rabbitmq"`
	}
	App struct {
		Name     string `yaml:"name"`
		Version  string `yaml:"version"`
		TimeZone int
		//EveryCodeMap map[int]int
	}
	//env-required:"true" -ОБЯЗАТЕЛЬНО должен получить переменную либо из окружения, либо из yaml.
	//TODO: Между true и false разницы не заметил. Разобраться

	Polycom struct {
		PolyUsername     string `env-required:"true" yaml:"poly_usernamename"    env:"POLY_USERNAME"`
		PolyPassword     string `env-required:"true" yaml:"poly_password"        env:"POLY_PASSWORD"`
		PolyEveryCodeMap map[int]bool
		RestartHour      int
		PolySwitch       int
		PolyMode         string
		PolyLogLevel     string
	}
	Lenovo struct {
		ZabbixUsername     string
		ZabbixPassword     string
		LenovoEveryCodeMap map[int]bool
		LenovoSwitch       int
		LenovoMode         string
		LenovoLogLevel     string
	}
	Ubiquiti struct {
		UiUsername      string `env-required:"true" yaml:"unifi_usernamename"   env:"UNIFI_USERNAME"`
		UiPassword      string `env-required:"true" yaml:"unifi_password"       env:"UNIFI_PASSWORD"`
		UiContrlRostov  string `env-required:"true" yaml:"contrl_rostov"   env:"UNIFI_CONTROLLER_ROSTOV"`
		UiContrlNovosib string `env-required:"true" yaml:"contrl_novosib"  env:"UNIFI_CONTROLLER_NOVOSIB"`
		Daily           int
		H1              int
		H2              int
		UiEveryCodeMap  map[int]int
		UiSwitch        int
		UiMode          string
		UiLogLevel      string
	}
	Eltex struct {
		EltexUsername     string
		EltexPassword     string
		EltexCntrlRostov  string
		EltexCntrlNovosib string
		EltexEveryCodeMap map[int]int
		EltexSwitch       int
		EltexMode         string
		EltexLogLevel     string
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
	DbGlpi struct {
		//строка подключения к серверу без указания БД
		GlpiConnectStr string `env-required:"true"   env:"GLPI_CONNECT_STR"`
		GlpiDB         string //задаю аргументами командной строки
	}
	DbGisupMySql struct {
		//строка подключения к серверу без указания БД
		//на текущий момент крутится на том же сервере, что и GLPI
		GisupConnectStr string `env-required:"true"   env:"GISUP_CONNECT_STR"`
		GisupDB         string //задаю аргументами командной строки
		//GisupDBprod     string //задаю аргументами командной строки
		//GisupDBtest     string //задаю аргументами командной строки
	}
	DbGisupPg struct {
		//строка подключения к серверу без указания БД
		PgConnectStr string `env:"PG_CONNECT_STR"`
		PgDb         string //задаю аргументами командной строки
		//PoolMax int `yaml:"pool_max" env:"PG_POOL_MAX"`
	}
	DbRedis struct {
		//строка подключения к серверу без указания БД
		RedisConnectString string `env:"REDIS_CONNECT_STR"`
		RedisDB            string //задаю аргументами командной строки
	}
	C3po struct {
		//Это универсальная строка для подключения и если будут добавляться новые методы, то не изменится
		//поэтому пока использования отдельно логина и пароля
		C3poUrl string `env-required:"true"   env:"C3PO_URL"`
		//C3poLogin    string `env-required:"true"   env:"C3PO_LOGIN"`
		//C3poPassword string `env-required:"true"   env:"C3PO_PASSWORD"`
	}
	Ldap struct {
		LdapDN       string `env-required:"true"   env:"LDAP_DN"`
		LdapDomain   string `env-required:"true"   env:"LDAP_Domain"`
		LdapLogin    string `env:"LDAP_LOGIN"`
		LdapPassword string `env:"LDAP_PASSWORD"`
		LdapRoleDn   string `env-required:"true"   env:"LDAP_ROLE_DN"`
		LdapServer   string `env-required:"true"   env:"LDAP_SERVER"`
	}

	Token struct {
		TTL    int
		JwtKey string `env-required:"true"   env:"JWT_KEY"`
	}
	HTTP struct {
		HttpURL  string `env-required:"true"   env:"GISUP_HTTP_URL"`
		HttpPort string //`yaml:"port" env:"HTTP_PORT"`
	}
	GRPC struct {
		GrpcPort int
	}
	RMQ struct {
		//ServerExchange 	string `yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
		RmqServerExchange string
		//ClientExchange string `yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
		RmqConnectStr string `env:"RMQ_CONNECT_STR"`
	}
	/*
		Log struct {
			LevelEnv string `yaml:"log_level"   env:"LOG_LEVEL"`
			LevelCmd string
		}*/
)
