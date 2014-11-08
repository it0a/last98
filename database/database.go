package database

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InitDB() {
	DB = NewDB()
}

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite3", "database/last98.sqlite")
	if err != nil {
		log.Fatal("Database initialization error: ", err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS images(id INTEGER PRIMARY KEY, description TEXT)")
	if err != nil {
		log.Fatal("Couldn't create books table!", err)
	}
	return db
}
