package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var DB *sql.DB

func InitDB() {
	DB = NewDB()
}

func NewDB() *sql.DB {
	db, err := sql.Open("postgres", "user=it0a dbname=last98 sslmode=disable")
	if err != nil {
		log.Fatal("Database initialization error: ", err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS images(id serial PRIMARY KEY, description TEXT)")
	if err != nil {
		log.Fatal("Couldn't create images table!", err)
	}
	return db
}
