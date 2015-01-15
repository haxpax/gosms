package gosms

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type SMSData struct {
	Messages             []SMS  `json:"messages"`
	IDisplayStart        string `json:"iDisplayStart"`
	IDisplayLength       int    `json:"iDisplayLength"`
	ITotalRecords        int    `json:"iTotalRecords"`
	ITotalDisplayRecords int    `json:"iTotalDisplayRecords"`
}

// Cache templates
var templates = template.Must(template.ParseFiles("gosmslib/templates/index.html", "gosmslib/templates/smsdata.html"))

/* web views */

//sms log viewer
func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- indexHandler")
	templates.ExecuteTemplate(w, "smsdata.html", nil)
}

//test sending sms
func testSMSHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println("--- testSMSHandlerGet")
	templates.ExecuteTemplate(w, "index.html", nil)
}

//handle test sms POST request
func testSMSHandlerPost(w http.ResponseWriter, r *http.Request) {
	log.Println("--- testSMSHandlerPost")
	r.ParseForm()
	mobile := r.FormValue("mobile")
	message := r.FormValue("message")
	SendSMS(mobile, message)
	w.Write([]byte("OK"))
}

//handle all static files based on specified path
//for now its /assets
func handleStatic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	static := vars["path"]
	http.ServeFile(w, r, filepath.Join("./gosmslib/assets", static))
}

/* end web views */

/* API views */

//push sms, allowed methods: POST
func smsAPIHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- smsAPIHandler")
	r.ParseForm()
	mobile := r.FormValue("mobile")
	message := r.FormValue("message")
	uuid := uuid.NewV1()
	sms := &SMS{UUID: uuid.String(), Mobile: mobile, Body: message}
	EnqueueMessage(sms)
	w.Write([]byte("OK"))
}

//dumps data, used by log view. Methods allowed: GET
func smsDataAPIHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- smsDataAPIHandler")
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

/* end API views */

func InitServer(host string, port string) error {
	log.Println("--- InitServer ", host, port)

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/", indexHandler)

	//handle static files
	r.HandleFunc(`/assets/{path:[a-zA-Z0-9=\-\/\.\_]+}`, handleStatic)

	//get rid of this view once done testing
	testsms := r.Path("/testsms/").Subrouter()
	testsms.Methods("GET").HandlerFunc(testSMSHandlerGet)
	testsms.Methods("POST").HandlerFunc(testSMSHandlerPost)

	//all API views
	r.HandleFunc("/api/sms/", smsAPIHandler)
	r.HandleFunc(`/api/smsdata/{start:[0-9]+}`, smsDataAPIHandler)

	http.Handle("/", r)

	bind := fmt.Sprintf("%s:%s", host, port)
	log.Println("listening on: ", bind)
	return http.ListenAndServe(bind, nil)

}
