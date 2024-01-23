package main

import (
	"github.com/deniskaponchik/GoSoft/config/ui"
	"github.com/deniskaponchik/GoSoft/internal/app"
	"io"
	"log"
	"os"
	"time"
)

func main() {

	//
	//STANDARD LOG
	fileNameUnifi := "Unifi_App_" + time.Now().Format("2006-01-02_15.04.05") + ".log"
	fileLogUnifi, err := os.OpenFile(fileNameUnifi, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	} else {
		multiWriter := io.MultiWriter(os.Stdout, fileLogUnifi)
		log.SetOutput(multiWriter)
	}
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime)) //удалить префикс времени в логах
	log.SetFlags(0)                                      //удалить префикс времени в логах
	log.Println("")

	//
	//CONFIG
	cfg, err := ui.NewConfigUnifi()
	//cfg, err := config.NewConfigPoly()
	//cfg, err := config.NewPolyConfig()
	if err == nil {
		log.Println("Конфиг создался успешно")
	} else {
		log.Fatalf("Config error: %s", err)
	}

	//
	//Zerro Log
	//zl := logger.New(cfg.Log.LevelCmd)
	//zl.Info("")

	//
	// app.Run(cfg)
	app.RunUnifi(cfg)
}
