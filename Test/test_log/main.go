package main

import (
	"log"
	"os"
)

func main() {

	file1, err := os.OpenFile("logs1.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file1.Close()
	file2, err := os.OpenFile("logs2.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()
	file3, err := os.OpenFile("logs3.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file3.Close()

	//InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	LoggerUnifi := log.New(file1, "", 0)
	LoggerWeb := log.New(file2, "", 0)
	LoggerPoly := log.New(file3, "", 0)

	//LoggerUnifi.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	//LoggerWeb.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	//LoggerPoly.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	LoggerUnifi.Println("Third start Unifi")
	LoggerWeb.Println("Third start Web")
	LoggerPoly.Println("Third start Poly")

	/*
		logfile1, err := os.Create("app1.log")
		if err != nil {
			log.Fatal(err)
		}
		defer logfile1.Close()
		log.SetOutput(logfile1)
	*/
}
