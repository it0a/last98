package images

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"github.com/gorilla/mux"
	"github.com/it0a/last98/database"
	"github.com/it0a/last98/page"
	"github.com/nfnt/resize"
	"html/template"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
)

type ImageEncoder interface {
	EncodedImage() (string, error)
	EncodedThumb() (string, error)
}

func (i ImageData) EncodedImage() string {
	encData := base64.StdEncoding.EncodeToString(i.Image)
	return encData
}

func (i ImageData) EncodedThumb() string {
	encData := base64.StdEncoding.EncodeToString(i.Thumb)
	return encData
}

type ImageRepository interface {
	FindById(id string) (ImageData, error)
	Delete(id string) error
	Save(n ImageData) error
}

type ImageDatabase struct{}

type ImageData struct {
	ID          int64
	Description string
	Image       []byte
	Thumb       []byte
}

func (i ImageDatabase) FindById(id string) (ImageData, error) {
	query := "SELECT id, description, data FROM images WHERE id = " + id
	result := database.DB.QueryRow(query)
	data := struct {
		ID          sql.NullInt64
		Description sql.NullString
		Image       []byte
	}{}
	err := result.Scan(&data.ID, &data.Description, &data.Image)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Unhandled error in GetImage:", err)
		return ImageData{}, err
	}
	image := ImageData{data.ID.Int64, data.Description.String, data.Image, nil}
	return image, nil
}

func (i ImageDatabase) Delete(id string) error {
	query := "DELETE FROM images WHERE id = $1"
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (i ImageDatabase) Save(id ImageData) error {
	stmt, err := database.DB.Prepare("INSERT INTO images(description, data, tn_data) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id.Description, id.Image, id.Thumb)
	if err != nil {
		log.Println("Failed to save new image data to the database: ", err)
		return err
	}
	return nil
}

func CreateThumbnail(imageData []byte) ([]byte, error) {
	log.Println("Creating thumbnail for new image...")
	reader := bytes.NewReader(imageData)
	image, err := jpeg.Decode(reader)
	if err != nil {
		log.Println("Failed to decode jpeg image data: ", err)
		return nil, err
	}
	tn := resize.Thumbnail(320, 320, image, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, tn, nil)
	if err != nil {
		log.Println("Couldn't encode the generated thumbnail!", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func SaveImage(id ImageData, ir ImageRepository) error {
	err := ir.Save(id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteImage(id string, ir ImageRepository) error {
	err := ir.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func IsEndOfRow(i int) bool {
	return ((i > 0) && (i%4 == 0))
}

func GetImages() []ImageData {
	query := "SELECT id, description, tn_data FROM images"
	result, err := database.DB.Query(query)
	if err != nil {
		log.Fatal("Error executing query: "+query, err)
	}
	images := []ImageData{}
	for result.Next() {
		data := struct {
			ID          sql.NullInt64
			Description sql.NullString
			Thumb       []byte
		}{}
		err := result.Scan(&data.ID, &data.Description, &data.Thumb)
		if err != nil {
			log.Fatal("ERROR!", err)
		}
		image := ImageData{data.ID.Int64, data.Description.String, nil, data.Thumb}
		if image.Description == "" {
			image.Description = "no description available"
		}
		images = append(images, image)
	}
	return images
}

func ImagesHandler(response http.ResponseWriter, request *http.Request) {
	data := struct {
		Page   page.Page
		Images []ImageData
	}{page.Page{"Images"}, GetImages()}
	funcs := template.FuncMap{"IsEndOfRow": IsEndOfRow}
	tmpl := make(map[string]*template.Template)
	tmpl["images.tmpl"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/base.tmpl", "templates/images.tmpl"))
	err := tmpl["images.tmpl"].ExecuteTemplate(response, "base", data)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func ReadImage(id string, imageRepository ImageRepository) (ImageData, error) {
	imageData, err := imageRepository.FindById(id)
	return imageData, err
}

func ImageShowHandler(response http.ResponseWriter, request *http.Request) {
	id := mux.Vars(request)["id"]
	image, err := ReadImage(id, ImageDatabase{})
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
	description := request.FormValue("description")
	file, _, err := request.FormFile("file")
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	thumb, err := CreateThumbnail(data)
	if err != nil {
		return
	}
	err = SaveImage(ImageData{Description: description, Image: data, Thumb: thumb}, ImageDatabase{})
	if err != nil {
		return
	}
	http.Redirect(response, request, "/images", http.StatusFound)
}

func ImagesDeleteHandler(response http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	DeleteImage(id, ImageDatabase{})
	http.Redirect(response, request, "/images", http.StatusFound)
}
