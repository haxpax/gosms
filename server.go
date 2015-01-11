package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Cache templates
var templates = template.Must(template.ParseFiles("index.html"))

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
	SendSMS(mobile, message)
	w.Write([]byte("OK"))
}

func InitServer(host string, port string) error {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/testsms/", testSMSHandler)

	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	bind := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("listening on %s...", bind)
	return http.ListenAndServe(bind, nil)

}
