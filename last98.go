package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Page struct {
	Title string
}

type Image struct {
	ID          sql.NullInt64
	Description sql.NullString
}

var db *sql.DB

func main() {
	log.Printf("Starting server...")
	start := time.Now()
	//
	var host = flag.String("host", "127.0.0.1", "IP of host to run web server on")
	var port = flag.Int("port", 8080, "Port to run webserver on")
	var staticPath = flag.String("staticPath", "static/", "Path to static files")
	flag.Parse()
	//
	db = NewDB()
	//
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/images", ImagesHandler)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(*staticPath))))
	//
	addr := fmt.Sprintf("%s:%d", *host, *port)
	elapsed := time.Since(start)
	log.Printf("Started in %s", elapsed)
	log.Printf("Listening on %s", addr)
	//
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

func IndexHandler(response http.ResponseWriter, request *http.Request) {
	log.Printf("Handling request with IndexHandler")
	page := Page{"Images"}
	data := struct {
		Page Page
	}{page}
	tmpl := make(map[string]*template.Template)
	tmpl["index.tmpl"] = template.Must(template.ParseFiles("templates/base.tmpl", "templates/index.tmpl"))
	err := tmpl["index.tmpl"].ExecuteTemplate(response, "base", data)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func getImages() []Image {
	query := "SELECT id, description FROM images"
	result, err := db.Query(query)
	if err != nil {
		log.Fatal("Error executing query: "+query, err)
	}
	log.Println("OK")
	images := []Image{}
	for result.Next() {
		image := Image{}
		err := result.Scan(&image.ID, &image.Description)
		if err != nil {
			log.Fatal("ERROR!", err)
		}
		images = append(images, image)
	}
	return images
}

func ImagesHandler(response http.ResponseWriter, request *http.Request) {
	log.Printf("Handling request with ImagesHandler")
	data := struct {
		Page   Page
		Images []Image
	}{Page{"Images"}, getImages()}
	tmpl := make(map[string]*template.Template)
	tmpl["images.tmpl"] = template.Must(template.ParseFiles("templates/base.tmpl", "templates/images.tmpl"))
	err := tmpl["images.tmpl"].ExecuteTemplate(response, "base", data)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite3", "last98.sqlite")
	if err != nil {
		log.Fatal("Database initialization error: ", err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS images(id INTEGER PRIMARY KEY, description TEXT)")
	if err != nil {
		log.Fatal("Couldn't create books table!", err)
	}
	return db
}
