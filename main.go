package main

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
	walletRepo := db.NewWalletRepo(sqlDB, conf)
	authService := services.NewAuthService(authRepo, walletRepo, conf)
	walletService := services.NewWalletService(walletRepo, conf)

	s := &server.Server{
		Config:           conf,
		AuthRepository:   authRepo,
		AuthService:      authService,
		WalletRepository: walletRepo,
		WalletService:    walletService,
	}
	//CRONJOBS
	go walletService.CronjobWebhookUpdate(walletService)
	go walletService.CronjobToReleaseLockedFunds()
	s.Start()
}
