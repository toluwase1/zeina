package zeina

import (
	"log"
	"net/http"
	"time"
	"zeina/config"
	"zeina/db"
	"zeina/server"
	"zeina/services"
)

func main() {
	http.DefaultClient.Timeout = time.Second * 10
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB := db.GetDB(conf)
	authRepo := db.NewAuthRepo(sqlDB)
	authService := services.NewAuthService(authRepo, conf)

	s := &server.Server{
		Config:         conf,
		AuthRepository: authRepo,
		AuthService:    authService,
	}
	s.Start()
}
