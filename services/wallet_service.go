package services

import (
	"log"
	"net/http"
	"zeina/config"
	"zeina/db"
	apiError "zeina/errors"
	"zeina/models"
)

type WalletService interface {
	CreateAccount(request *models.Account) (*models.Account, *apiError.Error)
}

// walletService struct
type walletService struct {
	Config     *config.Config
	walletRepo db.WalletRepository
}

// NewWalletService instantiate an walletService
func NewWalletService(walletRepo db.WalletRepository, conf *config.Config) WalletService {
	return &walletService{
		Config:     conf,
		walletRepo: walletRepo,
	}
}

func (a *walletService) CreateAccount(account *models.Account) (*models.Account, *apiError.Error) {
	account, err := a.walletRepo.CreateAccount(account)
	if err != nil {
		log.Printf("unable to create user: %v", err.Error())
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}

	return account, nil
}
