package main

import (
	"fmt"
	//"../config"
	//"../internal/app"
	"github.com/deniskaponchik/GoSoft/Unifi/config/ui"
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
	fmt.Println("")

	cfg, err := ui.NewConfigUnifi()
	//cfg, err := config.NewConfigPoly()
	//cfg, err := config.NewPolyConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// app.Run(cfg)
	app.RunUnifi(cfg)
}