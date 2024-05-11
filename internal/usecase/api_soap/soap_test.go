package api_soap

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"testing"
)

// cd internal\usecase\api_soap
// go test -run SoapGetDirectory
func Test_SoapGetDirectory(t *testing.T) {

	type (
		ConfigSoap struct {
			SoapProd string `env-required:"true" env:"SOAP_PROD"`
			SoapTest string `env-required:"true" env:"SOAP_TEST"`
		}
	)
	cfg := &ConfigSoap{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return //nil, err
	} else {
		log.Println(cfg.SoapProd)
		log.Println(cfg.SoapTest)
	}

	ss := NewSoap(cfg.SoapTest, "bpmUrl")

	mapRegions := make(map[string]bool)

	err = ss.GetDirectoryFieldsAll("UsrHelpDeskRegion", mapRegions)
	if err != nil {
		t.Errorf("Мапа регионов не была получена: %s", err)
	} else {
		log.Println("")
		//t.Logf("Мапа регионов загружена из bpm: ")
		log.Println("Мапа регионов загружена из bpm: ")
		for k, _ := range mapRegions {
			log.Println(k)
		}
	}
}

/*
func Test_InitializeSoapClient() (*Soap, error) {

	type(
		ConfigSoap struct {
			SoapProd string `env-required:"true" env:"SOAP_PROD"`
			SoapTest string `env-required:"true" env:"SOAP_TEST"`
		}
	)
	cfg := &ConfigSoap{	}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}


	soap := NewSoap(cfg.SoapTest, "bpm")
	return soap, nil
}
*/
