package main

import (
	"encoding/json"
	"fmt"
	"github.com/haxpax/gosms"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"encoding/base64"
)

//reposne structure to /sms
type SMSResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

//response structure to /smsdata/
type SMSDataResponse struct {
	Status   int            `json:"status"`
	Message  string         `json:"message"`
	Summary  []int          `json:"summary"`
	DayCount map[string]int `json:"daycount"`
	Messages []gosms.SMS    `json:"messages"`
}

// Cache templates
var templates = template.Must(template.ParseFiles("./templates/index.html"))

var authUsername string
var authPassword string

/* dashboard handlers */

// dashboard
func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- indexHandler")
	// templates.ExecuteTemplate(w, "index.html", nil)
	// Use during development to avoid having to restart server
	// after every change in HTML
	t, _ := template.ParseFiles("./templates/index.html")
	t.Execute(w, nil)
}

// handle all static files based on specified path
// for now its /assets
func handleStatic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	static := vars["path"]
	http.ServeFile(w, r, filepath.Join("./assets", static))
}

/* end dashboard handlers */

/* API handlers */

// push sms, allowed methods: POST
func sendSMSHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- sendSMSHandler")
	w.Header().Set("Content-type", "application/json")

	//TODO: validation
	r.ParseForm()
	mobile := r.FormValue("mobile")
	message := r.FormValue("message")
	uuid := uuid.NewV1()
	sms := &gosms.SMS{UUID: uuid.String(), Mobile: mobile, Body: message, Retries: 0}
	gosms.EnqueueMessage(sms, true)

	smsresp := SMSResponse{Status: 200, Message: "ok"}
	var toWrite []byte
	toWrite, err := json.Marshal(smsresp)
	if err != nil {
		log.Println(err)
		//lets just depend on the server to raise 500
	}
	w.Write(toWrite)
}

// dumps JSON data, used by log view. Methods allowed: GET
func getLogsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- getLogsHandler")
	messages, _ := gosms.GetMessages("")
	summary, _ := gosms.GetStatusSummary()
	dayCount, _ := gosms.GetLast7DaysMessageCount()
	logs := SMSDataResponse{
		Status:   200,
		Message:  "ok",
		Summary:  summary,
		DayCount: dayCount,
		Messages: messages,
	}
	var toWrite []byte
	toWrite, err := json.Marshal(logs)
	if err != nil {
		log.Println(err)
		//lets just depend on the server to raise 500
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(toWrite)
}

/* end API handlers */

func InitServer(host string, port string, username string, password string) error {
	log.Println("--- InitServer ", host, port)

	authUsername = username
	authPassword = password

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/", use(indexHandler, basicAuth))

	// handle static files
	r.HandleFunc(`/assets/{path:[a-zA-Z0-9=\-\/\.\_]+}`, use(handleStatic, basicAuth))

	// all API handlers
	api := r.PathPrefix("/api").Subrouter()
	api.Methods("GET").Path("/logs/").HandlerFunc(use(getLogsHandler, basicAuth))
	api.Methods("POST").Path("/sms/").HandlerFunc(use(sendSMSHandler, basicAuth))

	http.Handle("/", r)

	bind := fmt.Sprintf("%s:%s", host, port)
	log.Println("listening on: ", bind)
	return http.ListenAndServe(bind, nil)

}

// See https://gist.github.com/elithrar/7600878#comment-955958 for how to extend it to suit simple http.Handler's
func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(authUsername) == 0 {
			h.ServeHTTP(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 || pair[0] != authUsername || pair[1] != authPassword {
			http.Error(w, "Not authorized", 401)
			return
		}

		h.ServeHTTP(w, r)
	}
}