package main

import (
	"github.com/Omie/gosms/gosmslib"
	"log"
	"os"
)

func main() {

	log.Println("Initializing gosms")
	//load the config, abort if required config is not preset
	appConfig, err := gosms.GetConfig("conf.ini")
	if err != nil {
		log.Println("Invalid config: ", err.Error(), " Aborting")
		os.Exit(1)
	}

	db, err := gosms.InitDB("sqlite3", "db.sqlite")
	if err != nil {
		log.Println("Error initializing database: ", err, " Aborting")
		os.Exit(1)
	}
	defer db.Close()

	serverhost, _ := appConfig.Get("SETTINGS", "SERVERHOST")
	serverport, _ := appConfig.Get("SETTINGS", "SERVERPORT")
	comport, _ := appConfig.Get("DEVICE0", "COMPORT")

	err = gosms.InitModem(comport)
	if err != nil {
		log.Println("Error initializing modem: ", err.Error(), " Aborting")
		os.Exit(1)
	}

	log.Println("Initializing worker")
	gosms.InitWorker()

	log.Println("Initializing server")
	err = gosms.InitServer(serverhost, serverport)
	if err != nil {
		log.Println("Error starting server: ", err.Error(), " Aborting")
		os.Exit(1)
	}
}
