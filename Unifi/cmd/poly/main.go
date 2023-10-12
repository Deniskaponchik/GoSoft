package main

import (
	//"../config"
	//"../internal/app"
	"github.com/deniskaponchik/GoSoft/Unifi/config"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/app"
	"log"
)

/*
"../config"
"../internal/app"
"github.com/deniskaponchik/GoSoft/Unifi/config"
"github.com/deniskaponchik/GoSoft/Unifi/internal/app"
*/

func main() {
	cfg, err := config.NewConfig()
	//cfg, err := config.NewPolyConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// app.Run(cfg)
	app.PolyRun(cfg)
}
