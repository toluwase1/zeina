package models

import (
	"database/sql"
	"github.com/gobuffalo/pop/genny/config"
	"log"
)

func getPostgresDB(c *config.Config) *sql.DB {
	log.Printf("Connecting to postgres: %+v", c)
	connStr := "postgres://postgres:toluwase@localhost/zeina?sslmode=disable"

	// Open a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Ping the database to check if the connection is working
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
