package main

import (
	"fmt"
	"log"

	//"../config"
	//"../internal/app"
	"github.com/deniskaponchik/GoSoft/Unifi/config/ui"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/app"
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
	if err == nil {
		fmt.Println("Конфиг создался успешно")
	} else {
		log.Fatalf("Config error: %s", err)
	}

	// app.Run(cfg)
	app.RunUnifi(cfg)
}
