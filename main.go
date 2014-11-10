package main

import (
	"flag"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"last98/database"
	"last98/images"
	"last98/index"
	"log"
	"net/http"
	"os"
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
	err := http.ListenAndServe(":"+get_port(), router)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

func get_port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
