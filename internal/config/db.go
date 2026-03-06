package config

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDB() {

	dsn := "host=localhost port=5432 user=postgres password=123456 dbname=Hospital sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	DB = db
}
