package images

import (
	"database/sql"
	"log"
)

type Image struct {
	ID          sql.NullInt64
	Description sql.NullString
}

func GetImages(db *sql.DB) []Image {
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
