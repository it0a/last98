package database

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
	"os"
)

var DB *sql.DB

func InitDB() {
	DB = NewDB()
}

func NewDB() *sql.DB {
	url := os.Getenv("DATABASE_URL")
	log.Println("Attempting to connect with url => " + url)
	if url == "" {
		log.Fatal("DATABASE_URL is not set!")
	}
	connection, err := pq.ParseURL(url)
	if err != nil {
		log.Fatal("Failed parsing url => "+url, err)
	}
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal("Database initialization error: ", err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS images(id serial PRIMARY KEY, description TEXT, data BYTEA, tn_data BYTEA)")
	if err != nil {
		log.Fatal("Couldn't create images table!", err)
	}
	return db
}
