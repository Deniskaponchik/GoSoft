package gisup

type (
	ConfigGisup struct {
		Polycom
		Lenovo
		Ubiquiti
		Eltex
		Bpm
		Soap
		GLPI
		C3po
		Ldap

		App `yaml:"app"`
		Token
		HTTP `yaml:"http"`
		GRPC
		//Log `yaml:"logger"`
		PG  `yaml:"postgres"`
		RMQ `yaml:"rabbitmq"`
	}
	App struct {
		Name     string `yaml:"name"`
		Version  string `yaml:"version"`
		TimeZone int
		//EveryCodeMap map[int]int
	}
	//env-required:"true" -ОБЯЗАТЕЛЬНО должен получить переменную либо из окружения, либо из yaml.
	//Между true и false разницы не заметил. Разобраться

	Polycom struct {
		PolyUsername     string `env-required:"true" yaml:"poly_usernamename"    env:"POLY_USERNAME"`
		PolyPassword     string `env-required:"true" yaml:"poly_password"        env:"POLY_PASSWORD"`
		PolySwitch       int
		PolyLogLevel     string
		RestartHour      int
		PolyEveryCodeMap map[int]bool
	}
	Lenovo struct {
		ZabbixUsername     string
		ZabbixPassword     string
		LenovoSwitch       int
		LenovoLogLevel     string
		LenovoEveryCodeMap map[int]bool
	}
	Ubiquiti struct {
		UiUsername      string `env-required:"true" yaml:"unifi_usernamename"   env:"UNIFI_USERNAME"`
		UiPassword      string `env-required:"true" yaml:"unifi_password"       env:"UNIFI_PASSWORD"`
		UiContrlRostov  string `env-required:"true" yaml:"contrl_rostov"   env:"UNIFI_CONTROLLER_ROSTOV"`
		UiContrlNovosib string `env-required:"true" yaml:"contrl_novosib"  env:"UNIFI_CONTROLLER_NOVOSIB"`
		//UiContrlstr     string
		//UiContrlint     int //для совместного приложения двух контроллеров не должен приходить с конфигом
		Daily          int
		H1             int
		H2             int
		UiEveryCodeMap map[int]int
		UiSwitch       int
		UiLogLevel     string
	}
	Eltex struct {
		EltexUsername     string
		EltexPassword     string
		EltexCntrlRostov  string
		EltexCntrlNovosib string
		EltexEveryCodeMap map[int]int
		EltexSwitch       int
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
	GLPI struct {
		GlpiConnectStr string `env-required:"true"   env:"GLPI_CONNECT_STR"`
		//строка подключения к серверу без указания БД
		//GlpiConnectStrGLPI string `env-required:"true"   env:"GLPI_CONNECT_STR_GLPI"`
		//GlpiITsupportProd  string `env-required:"true"   env:"GLPI_CONNECT_STR_ITSUP"`
		//GlpiITsupportTest  string `env-required:"true"   env:"GLPI_ITSUP_TEST"`
		//GlpiITsupport      string //`env-required:"false"`
		DB string //имя базы данных для unifi таблиц. задаю аргументами командной строки
	}
	PG struct {
		//PoolMax int `yaml:"pool_max" env:"PG_POOL_MAX"`
		PgConnectStr string `env:"PG_CONNECT_STR"`
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
	/*
		Log struct {
			LevelEnv string `yaml:"log_level"   env:"LOG_LEVEL"`
			LevelCmd string
		}*/
	Token struct {
		TTL    int
		JwtKey string `env-required:"true"   env:"JWT_KEY"`
	}
	HTTP struct {
		URL  string `env-required:"true"   env:"GISUP_HTTP_URL"`
		Port string //`yaml:"port" env:"HTTP_PORT"`
	}
	GRPC struct {
		Port int
	}
	RMQ struct {
		//ServerExchange string `yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
		ServerExchange string
		//ClientExchange string `yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
		RmqConnectStr string `env:"RMQ_CONNECT_STR"`
	}
)
