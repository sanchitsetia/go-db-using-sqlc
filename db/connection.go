package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB
var queries *Queries

func InitDb() {
	var err error
	db, err = sql.Open("postgres", "postgresql://sanchit:password@localhost:5432/yourdatabase?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}
	queries = New(db)
}

func GetDB() (*sql.DB, *Queries) {
	return db, queries
}
