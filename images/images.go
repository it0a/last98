package images

import (
	"database/sql"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"last98/database"
	"last98/page"
	"log"
	"net/http"
	"os"
)

type Image struct {
	ID          sql.NullInt64
	Description sql.NullString
}

func SaveImage(request *http.Request) {
	file, header, err := request.FormFile("file")
	if err != nil {
		log.Fatal("Failed to retrieve file from form!", err)
	}
	defer file.Close()
	out, err := os.Create("images/data/" + header.Filename)
	if err != nil {
		log.Fatal("Failed to write file!", err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal("Failed to copy the POST request into a new file!", err)
	}
	log.Printf("Successfully saved file: %s", header.Filename)
	log.Printf("Writing record to the database...")
	description := request.FormValue("description")
	query := "INSERT INTO images(description) VALUES('" + description + "')"
	_, err = database.DB.Exec(query)
	if err != nil {
		log.Fatal("Couldn't save image!", err)
	}
}

func GetImages() []Image {
	query := "SELECT id, description FROM images"
	result, err := database.DB.Query(query)
	if err != nil {
		log.Fatal("Error executing query: "+query, err)
	}
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

func GetImage(id string) (Image, error) {
	query := "SELECT id, description FROM images WHERE id = " + id
	result := database.DB.QueryRow(query)
	image := Image{}
	err := result.Scan(&image.ID, &image.Description)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal("Unhandled error in GetImage:", err)
	}
	return image, err
}

func DeleteImage(id string) {
	log.Println("Deleting image ID " + id)
	query := "DELETE FROM images WHERE id = $1"
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		log.Fatal("Couldn't prepare deletion statement!", err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal("Couldn't delete image!", err)
	}
}

func ImagesHandler(response http.ResponseWriter, request *http.Request) {
	data := struct {
		Page   page.Page
		Images []Image
	}{page.Page{"Images"}, GetImages()}
	tmpl := make(map[string]*template.Template)
	tmpl["images.tmpl"] = template.Must(template.ParseFiles("templates/base.tmpl", "templates/images.tmpl"))
	err := tmpl["images.tmpl"].ExecuteTemplate(response, "base", data)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func ImageShowHandler(response http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"]
	image, err := GetImage(id)
	if err != nil {
		log.Println("Error occurred when retrieving image ID " + id + " - redirecting to images index")
		http.Redirect(response, request, "/images", http.StatusFound)
	}
	data := struct {
		Page  page.Page
		Image Image
	}{page.Page{"Images"}, image}
	tmpl := make(map[string]*template.Template)
	tmpl["image.tmpl"] = template.Must(template.ParseFiles("templates/base.tmpl", "templates/image.tmpl"))
	err = tmpl["image.tmpl"].ExecuteTemplate(response, "base", data)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func ImagesSaveHandler(response http.ResponseWriter, request *http.Request) {
	SaveImage(request)
	http.Redirect(response, request, "/images", http.StatusFound)
}

func ImagesDeleteHandler(response http.ResponseWriter, request *http.Request) {
	DeleteImage(request.FormValue("id"))
	http.Redirect(response, request, "/images", http.StatusFound)
}
