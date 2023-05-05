package db

import (
	"database/sql"
	_ "github.com/lib/pq"
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
	postgresDSN := "postgres://toluwase:toluwase@localhost/zeina?sslmode=disable"
	//postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d TimeZone=Africa/Lagos",
	//	c.PostgresHost, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresPort)
	log.Println(postgresDSN)
	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		log.Println("db connection error", err)
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Connected to database successfully!")
	return db
}
