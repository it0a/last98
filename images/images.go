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
	"strconv"
)

type ImageModel struct {
	ID          int64
	Description string
	Image       ImageData
	Thumb       ImageData
}

type ImageDataI interface {
	Encode() (string, error)
	CreateThumbnail([]byte, error)
}

type ImageData struct {
	Data []byte
}

func (i ImageData) Encode() string {
	return base64.StdEncoding.EncodeToString(i.Data)
}

func ReadImage(id int64, ir ImageRepository) (ImageModel, error) {
	imageModel, err := ir.FindById(id)
	if err != nil {
		log.Println("Failed to read image from database!: ", err)
	}
	return imageModel, err
}

type ImageRepository interface {
	FindById(id int64) (ImageModel, error)
	Delete(id int64) error
	Save(im ImageModel) error
}

type ImageDatabase struct{}

func (i ImageDatabase) FindById(id int64) (ImageModel, error) {
	query := "SELECT id, description, data FROM images WHERE id = " + strconv.FormatInt(id, 10)
	result := database.DB.QueryRow(query)
	data := struct {
		ID          sql.NullInt64
		Description sql.NullString
		Image       []byte
	}{}
	if err := result.Scan(&data.ID, &data.Description, &data.Image); err != nil && err != sql.ErrNoRows {
		return ImageModel{}, err
	}
	return ImageModel{ID: data.ID.Int64, Description: data.Description.String, Image: ImageData{data.Image}}, nil
}

func (i ImageDatabase) Delete(id int64) error {
	stmt, err := database.DB.Prepare("DELETE FROM images WHERE id = $1")
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(id); err != nil {
		return err
	}
	return nil
}

func (i ImageDatabase) Save(im ImageModel) error {
	stmt, err := database.DB.Prepare("INSERT INTO images(description, data, tn_data) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(im.Description, im.Image.Data, im.Thumb.Data); err != nil {
		return err
	}
	return nil
}

func SaveImage(im ImageModel, ir ImageRepository) error {
	if err := ir.Save(im); err != nil {
		log.Println("Failed to save new image data to the database: ", err)
		return err
	}
	return nil
}

func DeleteImage(id int64, ir ImageRepository) error {
	if err := ir.Delete(id); err != nil {
		log.Println("Failed to delete image data to the database: ", err)
		return err
	}
	return nil
}

func (i ImageModel) CreateThumbnail() ([]byte, error) {
	image, err := jpeg.Decode(bytes.NewReader(i.Image.Data))
	if err != nil {
		return nil, err
	}
	tn := resize.Thumbnail(320, 320, image, resize.Lanczos3)
	buf := new(bytes.Buffer)
	if err = jpeg.Encode(buf, tn, nil); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func IsEndOfRow(i int) bool {
	return ((i > 0) && (i%4 == 0))
}

func GetImages() ([]ImageModel, error) {
	query := "SELECT id, description, tn_data FROM images"
	result, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	images := []ImageModel{}
	for result.Next() {
		data := struct {
			ID          sql.NullInt64
			Description sql.NullString
			Thumb       []byte
		}{}
		if err := result.Scan(&data.ID, &data.Description, &data.Thumb); err != nil {
			return nil, err
		}
		image := ImageModel{ID: data.ID.Int64, Description: data.Description.String, Thumb: ImageData{data.Thumb}}
		if image.Description == "" {
			image.Description = "no description available"
		}
		images = append(images, image)
	}
	return images, nil
}

// Handlers

func ImagesHandler(response http.ResponseWriter, request *http.Request) {
	thumbs, err := GetImages()
	data := struct {
		Page   page.Page
		Images []ImageModel
	}{page.Page{"Images"}, thumbs}
	funcs := template.FuncMap{"IsEndOfRow": IsEndOfRow}
	tmpl := make(map[string]*template.Template)
	tmpl["images.tmpl"] = template.Must(template.New("").Funcs(funcs).ParseFiles("templates/base.tmpl", "templates/images.tmpl"))
	if err = tmpl["images.tmpl"].ExecuteTemplate(response, "base", data); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func ImageShowHandler(response http.ResponseWriter, request *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(request)["id"], 10, 64)
	if err != nil {
		log.Fatal("Can't handle int64 on this system - bailing out")
	}
	image, err := ReadImage(id, ImageDatabase{})
	if err != nil {
		log.Println("Error occurred when retrieving image ID - redirecting to images index")
		http.Redirect(response, request, "/images", http.StatusFound)
	}
	data := struct {
		Page  page.Page
		Image ImageModel
	}{page.Page{"Images"}, image}
	tmpl := make(map[string]*template.Template)
	tmpl["image.tmpl"] = template.Must(template.ParseFiles("templates/base.tmpl", "templates/image.tmpl"))
	if err = tmpl["image.tmpl"].ExecuteTemplate(response, "base", data); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func ImagesSaveHandler(response http.ResponseWriter, request *http.Request) {
	newImage := ImageModel{Description: request.FormValue("description")}
	file, _, err := request.FormFile("file")
	if err != nil {
		log.Println("Error occurred when retrieving the file from the request form.")
		return
	}
	iData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Error occurred when reading the uploaded image - redirecting to images index.")
		return
	}
	newImage.Image.Data = iData
	newImage.Thumb.Data, err = newImage.CreateThumbnail()
	if err != nil {
		log.Println("Error creating a thumbnail from this data - not saving an image.")
	} else {
		SaveImage(newImage, ImageDatabase{})
	}
	http.Redirect(response, request, "/images", http.StatusFound)
}

func ImagesDeleteHandler(response http.ResponseWriter, request *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(request)["id"], 10, 64)
	if err != nil {
		log.Fatal("Can't handle int64 on this system - bailing out")
	}
	DeleteImage(id, ImageDatabase{})
	http.Redirect(response, request, "/images", http.StatusFound)
}
