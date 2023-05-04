package db

import (
	"database/sql"
	"log"
	"zeina/config"
)

type SqlDB struct {
	DB *sql.DB
}

func GetDB(c *config.Config) *SqlDB {
	sqlDB := &SqlDB{}
	sqlDB.Init(c)
	return sqlDB
}

func (sql *SqlDB) Init(c *config.Config) {
	sql.DB = getPostgresDB(c)
}

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
	log.Println("Connected to database successfully!")
	return db
}
