package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/it0a/last98/database"
	"github.com/it0a/last98/images"
	"github.com/it0a/last98/index"
	"github.com/it0a/last98/initialize"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Printf("Starting last98...")
	start := time.Now()
	//
	var staticPath = flag.String("staticPath", "static/", "Path to static files")
	flag.Parse()
	//
	database.InitDB()
	//
	router := mux.NewRouter()
	router.HandleFunc("/", index.IndexHandler)
	//
	router.HandleFunc("/images", images.ImagesHandler).Methods("GET")
	router.HandleFunc("/images", images.ImagesSaveHandler).Methods("POST")
	router.HandleFunc("/images/delete", images.ImagesDeleteHandler).Methods("POST")
	imageRouter := router.PathPrefix("/images/{id}").Subrouter()
	imageRouter.Methods("GET").HandlerFunc(images.ImageShowHandler)
	//
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(*staticPath))))
	//
	elapsed := time.Since(start)
	log.Printf("Initialization finished in %s", elapsed)
	//
	var evr initialize.EnvVarReader
	log.Println("Using port " + initialize.ReadPort(evr))
	err := http.ListenAndServe(":"+initialize.ReadPort(evr), router)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
