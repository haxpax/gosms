package main

import (
	"fmt"
    "net/http"
	"os"
	"runtime"
)

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "hello, world from %s ", runtime.Version())
}

func main() {

    //load the config, abort if required config is not preset
    appConfig, err := getConfig("conf.ini")
    if err != nil {
        fmt.Println("Invalid config: ", err.Error())
        os.Exit(1)
    }

    serverhost, _ := appConfig.Get("SETTINGS", "SERVERHOST")
    serverport, _ := appConfig.Get("SETTINGS", "SERVERPORT")

	http.HandleFunc("/", hello)

	bind := fmt.Sprintf("%s:%s", serverhost, serverport)
	fmt.Printf("listening on %s...", bind)
	err = http.ListenAndServe(bind, nil)

	if err != nil {
		panic(err)
	}
}

