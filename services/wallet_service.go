package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
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
	Deposit(ctx context.Context, depositRequest models.TransactionRequest) error
	Withdrawal(ctx context.Context, depositRequest models.TransactionRequest) error
	ConfirmDepositFromWebhook(ctx context.Context, delta int64, _type string, account models.Account, reference string) error
	ConfirmWithdrawalFromWebhook(ctx context.Context, delta int64, _type string, account models.Account, reference string) error
	CronjobWebhookUpdate(service WalletService)
	LockBalance(ctx *gin.Context, locker models.LockFunds) error
	UnLockMaturedBalance() error
	CronjobToReleaseLockedFunds()
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

func (a *walletService) Deposit(ctx context.Context, depositRequest models.TransactionRequest) error {
	account, err := a.walletRepo.FindAccountByNumber(depositRequest.AccountNumber)
	if err != nil {
		return fmt.Errorf("account/user number does not exist %v %v", err, http.StatusBadRequest)
	}
	if account.AccountType != depositRequest.AccountType {
		return fmt.Errorf("(%s) account type specified does not exist: %v", depositRequest.AccountType, err)
	}
	if account.AccountNumber == a.Config.ZeinaAccountNumber {
		return fmt.Errorf("wrong account number %v", http.StatusBadRequest)
	}
	if depositRequest.Amount <= 0 {
		return fmt.Errorf("invalid amount %v", http.StatusBadRequest)
	}
	depositRequest.Reference = uuid.New().String()
	err = a.InternalMove(ctx, depositRequest.Amount, models.Deposit, *account, depositRequest.Reference)
	if err != nil {
		return err
	}
	depositRequest.Purpose = models.Deposit
	db.PushRequestToQueue(depositRequest)
	return nil
}

func (a *walletService) Withdrawal(ctx context.Context, withdrawalRequest models.TransactionRequest) error {
	account, err := a.walletRepo.FindAccountByNumber(withdrawalRequest.AccountNumber)
	if err != nil {
		return fmt.Errorf("account/user number does not exist %v %v", err, http.StatusBadRequest)
	}
	if account.AccountType != withdrawalRequest.AccountType {
		return fmt.Errorf("(%s) account type specified does not exist: %v", withdrawalRequest.AccountType, err)
	}
	if account.AccountNumber == a.Config.ZeinaAccountNumber {
		return fmt.Errorf("wrong account number %v", http.StatusBadRequest)
	}
	if withdrawalRequest.Amount <= 0 {
		return fmt.Errorf("invalid amount %v", http.StatusBadRequest)
	}
	if withdrawalRequest.Amount > account.AvailableBalance {
		return fmt.Errorf("insufficient balance %v", http.StatusBadRequest)
	}
	if account.AvailableBalance < withdrawalRequest.Amount {
		return fmt.Errorf("insufficient balance %v", http.StatusPaymentRequired)
	}

	withdrawalRequest.Reference = uuid.New().String()
	err = a.ExternalMove(ctx, withdrawalRequest.Amount, models.Withdrawal, *account, withdrawalRequest.Reference)
	if err != nil {
		return err
	}
	withdrawalRequest.Purpose = models.Withdrawal
	db.PushRequestToQueue(withdrawalRequest)
	return nil
}

func (a *walletService) InternalMove(ctx context.Context, delta int64, _type string, account models.Account, reference string) error {
	var (
		side2    = models.DebitEntry
		side     = models.CreditEntry
		txAmount = delta
		txStatus = models.StatusPending
	)

	if _type == models.Withdrawal {
		delta = -delta
		side = models.DebitEntry
		side2 = models.CreditEntry
		txStatus = models.StatusCompleted
	}

	timeCreated := time.Now().Unix()
	zeinaAccount, err := a.walletRepo.FindZeinaAccount(a.Config.ZeinaAccountNumber)
	if err != nil {
		return err
	}
	entries := []models.Entry{
		{
			AccName:   "Pending",
			AccountID: account.ID,
			Delta:     delta,
			Side:      side,
		},
		{
			AccountID: zeinaAccount.ID,
			Delta:     -delta,
			Side:      side2,
		},
	}

	ledger := models.Ledger{
		ID:          uuid.New().String(),
		CreatedAt:   timeCreated,
		Entries:     entries,
		AccountType: account.AccountType,
		Type:        _type,
	}

	transaction := models.Transaction{
		BaseModel: models.BaseModel{
			ID:        uuid.New().String(),
			CreatedAt: timeCreated,
			UpdatedAt: &timeCreated,
		},
		AccountID:        account.ID,
		Entry:            side,
		Purpose:          _type,
		Status:           txStatus,
		Change:           txAmount,
		AvailableBalance: account.AvailableBalance,
		PendingBalance:   account.PendingBalance + delta,
		Reference:        reference,
	}

	return a.walletRepo.InternalMove(ctx, ledger, transaction)
}

func (a *walletService) ExternalMove(ctx context.Context, delta int64, _type string, account models.Account, reference string) error {

	var (
		side2    = models.CreditEntry
		side     = models.DebitEntry
		txAmount = delta
		txStatus = models.StatusCompleted
	)

	if _type == models.Withdrawal {
		delta = -delta
		side = models.CreditEntry
		side2 = models.DebitEntry
		txStatus = models.StatusPending
	}

	timeCreated := time.Now().Unix()

	entries := []models.Entry{
		{
			AccName:   "pending", //add to ledger
			AccountID: account.ID,
			Delta:     -delta,
			Side:      side,
		},
		{
			AccName:   "available", ////add to ledger
			AccountID: account.ID,
			Delta:     delta,
			Side:      side2,
		},
	}

	ledger := models.Ledger{
		ID:          uuid.New().String(),
		CreatedAt:   timeCreated,
		Entries:     entries,
		AccountType: account.AccountType,
		Type:        _type,
	}

	transaction := models.Transaction{
		BaseModel: models.BaseModel{
			ID:        uuid.New().String(),
			CreatedAt: timeCreated,
			UpdatedAt: &timeCreated,
		},
		AccountID:        account.ID,
		Entry:            side2,
		Change:           txAmount,
		Purpose:          _type,
		Status:           txStatus,
		AvailableBalance: account.AvailableBalance + delta,
		PendingBalance:   account.PendingBalance - delta,
		Reference:        reference,
	}
	return a.walletRepo.ExternalMove(ctx, ledger, transaction)
}

func (a *walletService) ConfirmDepositFromWebhook(ctx context.Context, delta int64, _type string, account models.Account, reference string) error {
	return a.ExternalMove(ctx, delta, _type, account, reference)
}

func (a *walletService) ConfirmWithdrawalFromWebhook(ctx context.Context, delta int64, _type string, account models.Account, reference string) error {
	return a.InternalMove(ctx, delta, _type, account, reference)
}

func (a *walletService) LockBalance(ctx *gin.Context, locker models.LockFunds) error {
	account, err := a.walletRepo.FindAccountByNumber(locker.AccountNumber)
	if err != nil {
		return fmt.Errorf("account/user number does not exist %v %v", err, http.StatusBadRequest)
	}
	if account.AccountType != locker.AccountType {
		return fmt.Errorf("(%s) account type specified does not exist: %v", locker.AccountType, err)
	}
	if account.AccountNumber == a.Config.ZeinaAccountNumber {
		return fmt.Errorf("wrong account number %v", http.StatusBadRequest)
	}
	if locker.Amount <= 0 {
		return fmt.Errorf("invalid amount %v", http.StatusBadRequest)
	}
	if locker.Amount > account.AvailableBalance {
		return fmt.Errorf("insufficient balance %v", http.StatusBadRequest)
	}
	if account.AvailableBalance < locker.Amount {
		return fmt.Errorf("insufficient balance %v", http.StatusPaymentRequired)
	}
	return a.walletRepo.LockBalance(ctx, locker, *account)
}

func (a *walletService) UnLockMaturedBalance() error {
	return a.walletRepo.ReleaseDueFundsWhenDue()
}

func (a *walletService) CronjobWebhookUpdate(service WalletService) {
	func() {
		for {
			log.Println("CronjobWebhookUpdate cronjob running")
			for _, request := range db.GetAllRequestsFromQueue() {
				account, err := a.walletRepo.FindAccountByNumber(request.AccountNumber)
				if err != nil {
					return
				}
				if request.Purpose == models.Withdrawal {
					log.Println("withdrawal operation occurring", request.Reference)
					err = service.ConfirmWithdrawalFromWebhook(context.Background(), request.Amount, request.Purpose, *account, request.Reference)
					if err != nil {
						return
					}
				} else {
					log.Println("deposit operation occurring")
					log.Println("checking reference: ", request.Reference)
					err = service.ConfirmDepositFromWebhook(context.Background(), request.Amount, request.Purpose, *account, request.Reference)
					if err != nil {
						return
					}
				}
			}
			time.Sleep(15 * time.Second)
		}
	}()
	select {}
}

func (a *walletService) CronjobToReleaseLockedFunds() {
	func() {
		for {
			log.Println("release funds cronjob running")
			err := a.UnLockMaturedBalance()
			if err != nil {
				log.Println(err)
				return
			}
			time.Sleep(15 * time.Second)
		}
	}()
	select {}
}
