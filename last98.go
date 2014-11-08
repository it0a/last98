package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

type Page struct {
	Title string
}

func main() {
	var host = flag.String("host", "127.0.0.1", "IP of host to run web server on")
	var port = flag.Int("port", 8080, "Port to run webserver on")
	var staticPath = flag.String("staticPath", "static/", "Path to static files")
	flag.Parse()
	//
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(*staticPath))))
	//
	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Listening on %s", addr)
	//
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

func IndexHandler(response http.ResponseWriter, request *http.Request) {
	log.Printf("Handling request with IndexHandler")
	page := Page{"Index"}
	tmpl := make(map[string]*template.Template)
	tmpl["index.html"] = template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
	err := tmpl["index.html"].ExecuteTemplate(response, "base", page)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}
