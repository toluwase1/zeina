package zeina

import (
	"log"
	"net/http"
	"time"
	"zeina/config"
	"zeina/db"
)

func main() {
	http.DefaultClient.Timeout = time.Second * 10
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB := db.GetDB(conf)
	authRepo := db.NewAuthRepo(sqlDB)

}
