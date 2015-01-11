package main

import (
	"fmt"
	"github.com/tarm/goserial"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type GSMModem struct {
	port   string
	baud   int
	status bool
	conn   io.ReadWriteCloser
}

var modem *GSMModem

// Cache templates
var templates = template.Must(template.ParseFiles("index.html"))

func (m *GSMModem) Connect() error {
	c := &serial.Config{Name: m.port, Baud: m.baud}
	s, err := serial.OpenPort(c)
	if err == nil {
		m.status = true
		m.conn = s
	}
	return err
}

func (m *GSMModem) sendCommand(command string) {
	time.Sleep(time.Duration(500 * time.Millisecond))
	_, err := m.conn.Write([]byte(command + "\r"))
	if err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	// 404 for all other url path
	if r.URL.Path[1:] != "" {
		http.NotFound(w, r)
		return
	}

	templates.ExecuteTemplate(w, "index.html", nil)
}

func testSMSHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mobile := r.FormValue("mobile")
	message := r.FormValue("message")

	// Put Modem in SMS Text Mode
	modem.sendCommand("AT+CMGF=1")

	modem.sendCommand("AT+CMGS=\"" + mobile + "\"")
	modem.sendCommand(message)
	// EOM CTRL-Z
	modem.sendCommand(string(26))

	w.Write([]byte("OK"))
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
	comport, _ := appConfig.Get("SETTINGS", "COMPORT")

	modem = &GSMModem{port: comport, baud: 115200}
	err = modem.Connect()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/testsms/", testSMSHandler)

	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	bind := fmt.Sprintf("%s:%s", serverhost, serverport)
	fmt.Printf("listening on %s...", bind)
	err = http.ListenAndServe(bind, nil)

	if err != nil {
		panic(err)
	}
}
