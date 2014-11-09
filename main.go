package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"last98/database"
	"last98/images"
	"last98/index"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Printf("Starting server...")
	start := time.Now()
	//
	var host = flag.String("host", "127.0.0.1", "IP of host to run web server on")
	var port = flag.Int("port", 8080, "Port to run webserver on")
	var staticPath = flag.String("staticPath", "static/", "Path to static files")
	flag.Parse()
	//
	database.InitDB()
	//
	router := mux.NewRouter()
	router.HandleFunc("/", index.IndexHandler)
	router.HandleFunc("/images", images.ImagesHandler).Methods("GET")
	router.HandleFunc("/images", images.ImagesSaveHandler).Methods("POST")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(*staticPath))))
	//
	addr := fmt.Sprintf("%s:%d", *host, *port)
	elapsed := time.Since(start)
	log.Printf("Initialization finished in %s", elapsed)
	log.Printf("Listening on %s", addr)
	//
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
