package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"html/template"
	"net/http"
	"path/filepath"
)

// Cache templates
var templates = template.Must(template.ParseFiles("index.html", "smsdata.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// 404 for all other url path
	if r.URL.Path[1:] != "" {
		http.NotFound(w, r)
		return
	}
	templates.ExecuteTemplate(w, "smsdata.html", nil)
}

func testSMSHandlerGet(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func testSMSHandlerPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mobile := r.FormValue("mobile")
	message := r.FormValue("message")
	SendSMS(mobile, message)
	w.Write([]byte("OK"))
}

func SMSAPIHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mobile := r.FormValue("mobile")
	message := r.FormValue("message")
	uuid := uuid.NewV1()
	sms := &SMS{UUID: uuid.String(), Mobile: mobile, Body: message}
	EnqueueMessage(sms)
	w.Write([]byte("OK"))
}

type SMSData struct {
	Messages             []SMS  `json:"messages"`
	IDisplayStart        string `json:"iDisplayStart"`
	IDisplayLength       int    `json:"iDisplayLength"`
	ITotalRecords        int    `json:"iTotalRecords"`
	ITotalDisplayRecords int    `json:"iTotalDisplayRecords"`
}

func SMSDataAPIHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	startWith := vars["start"]
	filter := "LIMIT 10 OFFSET " + startWith
	messages, _ := getMessages(filter)
	smsdata := SMSData{
		Messages:             messages,
		IDisplayStart:        startWith,
		IDisplayLength:       10,
		ITotalRecords:        len(messages),
		ITotalDisplayRecords: len(messages),
	}
	var toWrite []byte
	toWrite, err := json.Marshal(smsdata)
	if err != nil {
		fmt.Println(err)
		toWrite = []byte("{\"messages\":[]}")
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(toWrite)
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	static := vars["path"]
	http.ServeFile(w, r, filepath.Join("./assets", static))
}

func InitServer(host string, port string) error {

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/", indexHandler)

	testsms := r.Path("/testsms/").Subrouter()
	testsms.Methods("GET").HandlerFunc(testSMSHandlerGet)
	testsms.Methods("POST").HandlerFunc(testSMSHandlerPost)

	r.HandleFunc("/api/sms/", SMSAPIHandler)
	r.HandleFunc(`/api/smsdata/{start:[0-9]+}`, SMSDataAPIHandler)
	r.HandleFunc(`/assets/{path:[a-zA-Z0-9=\-\/\.\_]+}`, handleStatic)

	http.Handle("/", r)

	bind := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("listening on %s...", bind)
	return http.ListenAndServe(bind, nil)

}
