package main

import (
	"io"
	"log"
	"os"
	"time"
)

func main() {

	fileNameUnifi := "Unifi_App_" + time.Now().Format("2006-01-02_15.04.05") + ".log"
	fileLogUnifi, err := os.OpenFile(fileNameUnifi, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	} else {
		multiWriter := io.MultiWriter(os.Stdout, fileLogUnifi)
		log.SetOutput(multiWriter)
	}

	log.Println("Success")
}
