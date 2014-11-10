package images

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"html/template"
	"image/jpeg"
	"io/ioutil"
	"last98/database"
	"last98/page"
	"log"
	"net/http"
)

type ImageData struct {
	ID          sql.NullInt64
	Description sql.NullString
	Data        string
}

func CreateThumbnail(imageData []byte) []byte {
	reader := bytes.NewReader(imageData)
	image, err := jpeg.Decode(reader)
	if err != nil {
		log.Fatal("Failed to generate thumbnail!", err)
	}
	tn := resize.Thumbnail(320, 320, image, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, tn, nil)
	if err != nil {
		log.Fatal("Couldn't encode the generated thumbnail!", err)
	}
	return buf.Bytes()
}

func SaveImage(request *http.Request) {
	log.Println("Reading request data...")
	description := request.FormValue("description")
	file, _, err := request.FormFile("file")
	if err != nil {
		log.Fatal("Failed to retrieve image from request!", err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Couldn't read the file into a byte array...", err)
	}
	//
	log.Println("Creating thumbnail for new image...")
	tn_data := CreateThumbnail(data)
	log.Println("Writing new image data...")
	stmt, err := database.DB.Prepare("INSERT INTO images(description, data, tn_data) VALUES($1, $2, $3)")
	if err != nil {
		log.Fatal("Error preparing insert statement!", err)
	}
	_, err = stmt.Exec(description, data, tn_data)
	if err != nil {
		log.Fatal("Couldn't save new image data!", err)
	}
}

func GetImages() []ImageData {
	query := "SELECT id, description, tn_data FROM images"
	result, err := database.DB.Query(query)
	if err != nil {
		log.Fatal("Error executing query: "+query, err)
	}
	images := []ImageData{}
	for result.Next() {
		image := ImageData{}
		data := []byte{}
		err := result.Scan(&image.ID, &image.Description, &data)
		if err != nil {
			log.Fatal("ERROR!", err)
		}
		image.Data = base64.StdEncoding.EncodeToString(data)
		if image.Description.String == "" {
			image.Description.String = "no description available"
		}
		images = append(images, image)
	}
	return images
}

func GetImage(id string) (ImageData, error) {
	query := "SELECT id, description, data FROM images WHERE id = " + id
	result := database.DB.QueryRow(query)
	image := ImageData{}
	data := []byte{}
	err := result.Scan(&image.ID, &image.Description, &data)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal("Unhandled error in GetImage:", err)
	}
	image.Data = base64.StdEncoding.EncodeToString(data)
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
		Images []ImageData
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
		Image ImageData
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
