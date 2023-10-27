package main

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	/*
		every66Code := map[int]bool{
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
		every63Code := map[int]bool{ //6 minutes
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
	*/
	every20Code := map[int]bool{
		5:  true,
		25: true,
		45: true,
	}

	cfg.BpmUrl = cfg.BpmTest
	cfg.SoapUrl = cfg.SoapTest

	cfg.EveryCode = every20Code
	cfg.Count20minute = 0
	cfg.CountHourFromDB = 0
	cfg.CountHourToDB = 0
	cfg.Reboot = 0

	return cfg, nil
}

type (
	Config struct {
		TempVars

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

	TempVars struct {
		EveryCode       map[int]bool `env-required:"false"`
		Count20minute   int          `env-default:"true"`
		CountHourFromDB int          `env-default:"true"`
		CountHourToDB   int          `env-default:"true"`
		Reboot          int          `env-default:"true"`
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
		GlpiConnectStrITsupport string `env-required:"true"   env:"GLPI_CONNECT_STR_ITSUP"`
		GlpiConnectStrGLPI      string `env-required:"true"   env:"GLPI_CONNECT_STR_GLPI"`
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
