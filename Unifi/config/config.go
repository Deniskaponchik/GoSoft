package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

func NewConfig(mode string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	if mode == "TEST" {
		cfg.BpmUrl = cfg.BpmTest
		cfg.SoapUrl = cfg.SoapTest
		cfg.GlpiITsupport = cfg.GlpiITsupportTest
	} else {
		cfg.BpmUrl = cfg.BpmProd
		cfg.SoapUrl = cfg.SoapProd
		cfg.GlpiITsupport = cfg.GlpiITsupportProd
	}

	return cfg, nil
}

type (
	Config struct {
		Polycom
		Ubiquiti
		Bpm
		Soap
		GLPI

		App  `yaml:"app"`
		HTTP `yaml:"netdial"`
		Log  `yaml:"logger"`
		PG   `yaml:"postgres"`
		RMQ  `yaml:"rabbitmq"`
	}
	Mode struct {
	}

	Polycom struct {
		PolyUsername string `env-required:"true" yaml:"poly_usernamename"    env:"POLY_USERNAME"`
		PolyPassword string `env-required:"true" yaml:"poly_password"        env:"POLY_PASSWORD"`
	}
	Ubiquiti struct {
		UiUsername      string `env-required:"true" yaml:"unifi_usernamename"   env:"UNIFI_USERNAME"`
		UiPassword      string `env-required:"true" yaml:"unifi_password"       env:"UNIFI_PASSWORD"`
		UiContrlRostov  string `env-required:"true" yaml:"contrl_rostov"   env:"UNIFI_CONTROLLER_ROSTOV"`
		UiContrlNovosib string `env-required:"true" yaml:"contrl_novosib"  env:"UNIFI_CONTROLLER_NOVOSIB"`
	}
	Bpm struct {
		BpmUrl  string `env-required:"false"`
		BpmProd string `env-required:"true" yaml:"bpm_prod"   env:"BPM_PROD"`
		BpmTest string `env-required:"true" yaml:"bpm_test"   env:"BPM_TEST"`
	}
	Soap struct {
		SoapUrl  string `env-required:"false"`
		SoapProd string `env-required:"true" env:"SOAP_PROD"`
		SoapTest string `env-required:"true" env:"SOAP_TEST"`
	}
	GLPI struct {
		GlpiConnectStrGLPI string `env-required:"true"   env:"GLPI_CONNECT_STR_GLPI"`
		GlpiITsupportProd  string `env-required:"true"   env:"GLPI_CONNECT_STR_ITSUP"`
		GlpiITsupportTest  string `env-required:"true"   env:"GLPI_ITSUP_TEST"`
		GlpiITsupport      string `env-required:"false"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
	Log struct {
		Level string `yaml:"log_level"   env:"LOG_LEVEL"`
	}
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}
	PG struct {
		PoolMax int `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		//URL     string `env-required:"true"                 env:"PG_URL"`
	}
	RMQ struct {
		ServerExchange string `env-required:"true" yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
		ClientExchange string `env-required:"true" yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
		//URL            string `env-required:"true"                            env:"RMQ_URL"`
	}
)
