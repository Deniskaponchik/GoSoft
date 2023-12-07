package main

import (
	"io"
	"log"
	"os"
	"time"

	//"../config"
	//"../internal/app"
	"github.com/deniskaponchik/GoSoft/Unifi/config/ui"
	"github.com/deniskaponchik/GoSoft/Unifi/internal/app"
)

func main() {
	log.Println("")

	//https://stackoverflow.com/questions/48629988/remove-timestamp-prefix-from-go-logger
	//time.Now().Format("2006-01-02 15:04:05")
	FileNameUnifi := "log_Unifi_App_" + time.Now().Format("2006-01-02_15.04.05") + ".log"
	fileLogUnifi, err := os.OpenFile(FileNameUnifi, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	multiWriter := io.MultiWriter(os.Stdout, fileLogUnifi)
	log.SetOutput(multiWriter)                           //(fileLogUnifi)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime)) //удалить префикс времени в логах
	log.SetFlags(0)                                      //удалить префикс времени в логах
	log.Println("")
	//Zerro Log  	//l := logger.New(cfg.Log.Level)   //l.Info("")

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
	// app.Run(cfg)
	app.RunUnifi(cfg)
}
