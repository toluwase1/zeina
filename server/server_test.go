package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
	"zeina/config"
	"zeina/db"
	"zeina/services"
)

var (
	server *Server
)

var testServer struct {
	router  *gin.Engine
	handler *Server
}

func TestMain(m *testing.M) {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("couldn't load env vars: %v", err)
	}
	fmt.Println("Starting server tests")
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB := db.GetDB(conf)
	authRepo := db.NewAuthRepo(sqlDB)
	walletRepo := db.NewWalletRepo(sqlDB, conf)
	authService := services.NewAuthService(authRepo, walletRepo, conf)
	walletService := services.NewWalletService(walletRepo, conf)
	testServer.handler = &Server{
		Config:           conf,
		AuthRepository:   authRepo,
		AuthService:      authService,
		WalletRepository: walletRepo,
		WalletService:    walletService,
	}
	testServer.handler.Config.JWTSecret = "testSecret"
	testServer.router = testServer.handler.setupRouter()
	exitCode := m.Run()
	os.Exit(exitCode)
}
