package main

import (
	"fmt"
	"github.com/haxpax/gosms"
	"github.com/haxpax/gosms/modem"
	"log"
	"os"
	"strconv"
)

func main() {

	log.Println("main: ", "Initializing gosms")
	//load the config, abort if required config is not preset
	appConfig, err := gosms.GetConfig("conf.ini")
	if err != nil {
		log.Println("main: ", "Invalid config: ", err.Error(), " Aborting")
		os.Exit(1)
	}

	db, err := gosms.InitDB("sqlite3", "db.sqlite")
	if err != nil {
		log.Println("main: ", "Error initializing database: ", err, " Aborting")
		os.Exit(1)
	}
	defer db.Close()

	serverhost, _ := appConfig.Get("SETTINGS", "SERVERHOST")
	serverport, _ := appConfig.Get("SETTINGS", "SERVERPORT")

	serverusername, _ := appConfig.Get("SETTINGS", "USERNAME")
	serverpassword, _ := appConfig.Get("SETTINGS", "PASSWORD")

	_numDevices, _ := appConfig.Get("SETTINGS", "DEVICES")
	numDevices, _ := strconv.Atoi(_numDevices)
	log.Println("main: number of devices: ", numDevices)

	var modems []*modem.GSMModem
	for i := 0; i < numDevices; i++ {
		dev := fmt.Sprintf("DEVICE%v", i)
		_port, _ := appConfig.Get(dev, "COMPORT")
		_baud := 115200 //appConfig.Get(dev, "BAUDRATE")
		_devid, _ := appConfig.Get(dev, "DEVID")
		m := modem.New(_port, _baud, _devid)
		modems = append(modems, m)
	}

	_bufferSize, _ := appConfig.Get("SETTINGS", "BUFFERSIZE")
	bufferSize, _ := strconv.Atoi(_bufferSize)

	_bufferLow, _ := appConfig.Get("SETTINGS", "BUFFERLOW")
	bufferLow, _ := strconv.Atoi(_bufferLow)

	_loaderTimeout, _ := appConfig.Get("SETTINGS", "MSGTIMEOUT")
	loaderTimeout, _ := strconv.Atoi(_loaderTimeout)

	_loaderCountout, _ := appConfig.Get("SETTINGS", "MSGCOUNTOUT")
	loaderCountout, _ := strconv.Atoi(_loaderCountout)

	_loaderTimeoutLong, _ := appConfig.Get("SETTINGS", "MSGTIMEOUTLONG")
	loaderTimeoutLong, _ := strconv.Atoi(_loaderTimeoutLong)

	log.Println("main: Initializing worker")
	gosms.InitWorker(modems, bufferSize, bufferLow, loaderTimeout, loaderCountout, loaderTimeoutLong)

	log.Println("main: Initializing server")
	err = InitServer(serverhost, serverport, serverusername, serverpassword)
	if err != nil {
		log.Println("main: ", "Error starting server: ", err.Error(), " Aborting")
		os.Exit(1)
	}
}
