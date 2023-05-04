package zeina

import (
	"database/sql"
	"fmt"
)

func main() {
	// Connection string for the PostgreSQL database
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

	fmt.Println("Connected to database successfully!")
}
