package main

import (
	//b "Unifi/config"
	// "./config"
	//"../GoSoft/Unifi/config"
	"log"
)

func main() {
	//cfg, err := config.NewConfig()
	cfg, err := config.NewPolyConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
