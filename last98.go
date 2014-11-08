package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(false)
	do_route(router)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", router)
}

func do_route(router *mux.Router) {
	do_images_route(router)
}

func do_images_route(router *mux.Router) {
	images := router.Path("/").Subrouter()
	images.Methods("GET").HandlerFunc(imagesIndexHandler)
	images.Methods("POST").HandlerFunc(imagesCreateHandler)
	do_images_singular_route(router)
}

func do_images_singular_route(router *mux.Router) {
	image := router.PathPrefix("/{id}").Subrouter()
	image.Methods("GET").Path("/edit").HandlerFunc(imageEditHandler)
	image.Methods("GET").HandlerFunc(imageShowHandler)
	image.Methods("PUT", "POST").HandlerFunc(imageUpdateHandler)
	image.Methods("DELETE").HandlerFunc(imageDeleteHandler)
}

func HomeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Home")
}

func imagesIndexHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "images index")
}

func imagesCreateHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "images create")
}

func imageShowHandler(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintln(rw, "showing image", id)
}

func imageUpdateHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "image update")
}

func imageDeleteHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "image delete")
}

func imageEditHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "image edit")
}
