package main

import (
	"fmt"
	"os"
)

func main() {

	//load the config, abort if required config is not preset
	appConfig, err := getConfig("conf.ini")
	if err != nil {
		fmt.Println("Invalid config: ", err.Error())
		os.Exit(1)
	}

	err = initDB("sqlite3", "db.sqlite")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	serverhost, _ := appConfig.Get("SETTINGS", "SERVERHOST")
	serverport, _ := appConfig.Get("SETTINGS", "SERVERPORT")
	comport, _ := appConfig.Get("SETTINGS", "COMPORT")

	err = InitModem(comport)
	if err != nil {
		fmt.Println("Error opening port: ", err.Error())
		os.Exit(1)
	}

	InitWorker()

	err = InitServer(serverhost, serverport)
	if err != nil {
		panic(err)
	}
}
