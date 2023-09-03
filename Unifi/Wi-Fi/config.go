package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type WiFiConfig struct {
	//OneDrive string
	UnifiUsername              string
	UnifiPassword              string
	UnifiControllerRostov      string
	UnifiControllerNovosib     string
	BpmProd                    string
	BpmTest                    string
	SoapProd                   string
	SoapTest                   string
	GlpiConnectStringITsupport string
	GlpiConnectStringGlpi      string
}

type WiFiConfigExt struct {
	WiFi WiFiConfig
	//DebugMode bool
	//
	//UserRoles []string
	Paths []string
	//MaxUsers  int
	ZesEnableSysman int
}

func NewWiFiConfig() *WiFiConfig {
	return &WiFiConfig{
		//OneDrive: getEnv("OneDrive", ""),
		UnifiUsername:              getEnv("UNIFI_USERNAME", ""),
		UnifiPassword:              getEnv("UNIFI_PASSWORD", ""),
		UnifiControllerRostov:      getEnv("UNIFI_CONTROLLER_ROSTOV", ""),
		UnifiControllerNovosib:     getEnv("UNIFI_CONTROLLER_NOVOSIB", ""),
		BpmProd:                    getEnv("BPM_PROD", ""),
		BpmTest:                    getEnv("BPM_TEST", ""),
		SoapProd:                   getEnv("SOAP_PROD", ""),
		SoapTest:                   getEnv("SOAP_TEST", ""),
		GlpiConnectStringITsupport: getEnv("GLPI_CONNECT_STR_ITSUP", ""),
		GlpiConnectStringGlpi:      getEnv("GLPI_CONNECT_STR_GLPI", ""),
	}
}

// New returns a new Config struct
func NewWiFiConfigExt() *WiFiConfigExt {
	return &WiFiConfigExt{
		WiFi: WiFiConfig{
			//OneDrive: getEnv("OneDrive", ""),
			/*
				UnifiUsername:              getEnv("UNIFI_USERNAME", ""),
				UnifiPassword:              getEnv("UNIFI_PASSWORD", ""),
				UnifiControllerRostov:      getEnv("UNIFI_CONTROLLER_ROSTOV", ""),
				UnifiControllerNovosib:     getEnv("UNIFI_CONTROLLER_NOVOSIB", ""),
				BpmProd:                    getEnv("BPM_PROD", ""),
				BpmTest:                    getEnv("BPM_TEST", ""),
				SoapProd:                   getEnv("SOAP_PROD", ""),
				SoapTest:                   getEnv("SOAP_TEST", ""),
				GlpiConnectStringITsupport: getEnv("GLPI_CONNECT_STR_ITSUP", ""),
				GlpiConnectStringGlpi:      getEnv("GLPI_CONNECT_STR_GLPI", ""),
			*/
		},
		//DebugMode: getEnvAsBool("DEBUG_MODE", true),
		//
		//UserRoles: getEnvAsSlice("USER_ROLES", []string{"admin"}, ","),
		Paths: getEnvAsSlice("Path", []string{"admin"}, ";"),
		//MaxUsers:  getEnvAsInt("MAX_USERS", 1),
		ZesEnableSysman: getEnvAsInt("ZES_ENABLE_SYSMAN", 1),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	} else {
		log.Fatalln("Не удалось получить переменную окружения: " + key)
	}
	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}
	val := strings.Split(valStr, sep)

	return val
}
