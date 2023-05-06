package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
	"zeina/config"
	"zeina/db"
	apiError "zeina/errors"
	"zeina/models"
)

type WalletService interface {
	CreateAccount(request *models.Account) (*models.Account, *apiError.Error)
	Deposit(ctx context.Context, user *models.User, depositRequest models.DepositRequest) error
	//CreateZeinaLedgerRecord(l *models.Ledger, zeinaAccount models.Account) error
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

func (a *walletService) Deposit(ctx context.Context, user *models.User, depositRequest models.DepositRequest) error {
	account, err := a.walletRepo.FindAccountByNumber("8163608141")
	if err != nil {
		return fmt.Errorf("account/user number does not exist %v %v", err, http.StatusBadRequest)
	}
	if account.AccountType != depositRequest.AccountType {
		return fmt.Errorf("(%s) account type specified does not exist: %v", depositRequest.AccountType, err)
	}
	if account.AccountNumber == a.Config.ZeinaAccountNumber {
		return fmt.Errorf("wrong account number %v", http.StatusBadRequest)
	}
	//zeinaAccount, err := a.walletRepo.FindZeinaAccount(a.Config.ZeinaAccountNumber)
	//if err != nil {
	//	return fmt.Errorf("internal server error %v", http.StatusInternalServerError)
	//}
	if depositRequest.Amount <= 0 {
		return fmt.Errorf("invalid amount %v", http.StatusBadRequest)
	}

	timeCreated := time.Now().Unix()
	ledger := models.Ledger{
		ID:          uuid.New().String(),
		CreatedAt:   timeCreated,
		AccountID:   account.ID,
		Entry:       models.CreditEntry,
		Change:      depositRequest.Amount,
		AccountType: depositRequest.AccountType,
		Type:        models.Deposit,
	}

}

func (a *walletService) InternalMove(ctx context.Context, ledger models.Ledger, account models.Account) error {
	if ledger.Entry == models.DebitEntry {
		ledger.Change = -ledger.Change
		ledger.Entry = models.CreditEntry
	}
	txTime := time.Now().Unix()
	account.PendingBalance += ledger.Change
	transaction := models.Transaction{
		BaseModel: models.BaseModel{
			ID:        uuid.New().String(),
			CreatedAt: txTime,
			UpdatedAt: &txTime,
		},
		AccountID:        account.ID,
		Entry:            models.CreditEntry,
		Purpose:          models.Deposit,
		Description:      models.Deposit,
		Remark:           models.Deposit,
		TotalBalance:     0,
		AvailableBalance: 0,
		PendingBalance:   account.PendingBalance,
	}

	return a.walletRepo.InternalMove(ctx, ledger, transaction)
}
