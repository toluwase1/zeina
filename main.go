package zeina

import (
	"log"
	"net/http"
	"time"
	"zeina/config"
	"zeina/models"
)

func main() {
	http.DefaultClient.Timeout = time.Second * 10
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	gormDB := models.GetDB(conf)

}
