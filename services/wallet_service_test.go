package services

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
	"zeina/config"
	"zeina/db"
	"zeina/models"
)

func TestDepositTransaction(t *testing.T) {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("couldn't load env vars: %v", err)
	}
	fmt.Println("Starting server tests")
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("...................", "")
	sqlDB := db.GetDB(conf)
	walletRepo := db.NewWalletRepo(sqlDB, conf)
	authRepo := db.NewAuthRepo(sqlDB)
	walletservice := NewWalletService(walletRepo, conf)
	_, account, _ := createAdminUser(t, walletRepo, authRepo)

	depositReq := models.TransactionRequest{
		Amount:        1000,
		AccountType:   models.AccountTypeSavings,
		AccountNumber: account.AccountNumber,
		Purpose:       models.Deposit,
		Reference:     uuid.New().String(),
	}
	require.NoError(t, walletservice.Deposit(context.Background(), depositReq))
	require.NoError(t, walletservice.Deposit(context.Background(), depositReq))
	require.NoError(t, walletservice.Deposit(context.Background(), depositReq))
	require.NoError(t, walletservice.Deposit(context.Background(), depositReq))
	require.NoError(t, walletservice.Deposit(context.Background(), depositReq))
	balance, err := walletservice.GetBalance(account.AccountNumber)
	require.NoError(t, err)
	require.Equal(t, int64(5000), balance.PendingBalance)

	//validate that ledger change column always sums to zero
	sum, err := walletservice.ValidateCorrectnessOfLedgersTable()
	require.NoError(t, err)
	require.Equal(t, int64(0), sum)

}

func TestWithdrawalTransaction(t *testing.T) {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("couldn't load env vars: %v", err)
	}
	fmt.Println("Starting server tests")
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("...................", "")
	sqlDB := db.GetDB(conf)
	walletRepo := db.NewWalletRepo(sqlDB, conf)
	authRepo := db.NewAuthRepo(sqlDB)
	walletservice := NewWalletService(walletRepo, conf)
	_, account, _ := createAdminUser(t, walletRepo, authRepo)

	depositReq := models.TransactionRequest{
		Amount:        1000,
		AccountType:   models.AccountTypeSavings,
		AccountNumber: account.AccountNumber,
		Purpose:       models.Deposit,
		Reference:     uuid.New().String(),
	}
	//test for insufficient  balance
	require.Error(t, walletservice.Withdrawal(context.Background(), depositReq))

	require.NoError(t, walletservice.Deposit(context.Background(), depositReq))
	require.NoError(t, walletservice.Deposit(context.Background(), depositReq))
	require.NoError(t, walletservice.Deposit(context.Background(), depositReq))
	require.NoError(t, walletservice.Deposit(context.Background(), depositReq))
	require.NoError(t, walletservice.Deposit(context.Background(), depositReq))

	//validate balance
	balance, err := walletservice.GetBalance(account.AccountNumber)
	require.NoError(t, err)
	require.Equal(t, int64(5000), balance.PendingBalance)

	require.NoError(t, walletservice.Withdrawal(context.Background(), depositReq))
	balance, err = walletservice.GetBalance(account.AccountNumber)
	require.NoError(t, err)
	require.Equal(t, int64(5000), balance.PendingBalance)

	//validate that ledger change column always sums to zero
	sum, err := walletservice.ValidateCorrectnessOfLedgersTable()
	require.NoError(t, err)
	require.Equal(t, int64(0), sum)

}

func createAdminUser(t *testing.T, walletRepo db.WalletRepository, authRepo db.AuthRepository) (*models.User, *models.Account, error) {
	t.Helper()
	accountTimeCreated := time.Now().Unix()
	password := gofakeit.Password(true, true, true, true, false, 10)
	hashPassword, err := GenerateHashPassword(password)
	if err != nil {
		return nil, nil, err
	}
	user := &models.User{
		BaseModel: models.BaseModel{
			ID:        uuid.New().String(),
			CreatedAt: accountTimeCreated,
			UpdatedAt: &accountTimeCreated,
		},
		Name:           gofakeit.Username(),
		HashedPassword: hashPassword,
		Email:          gofakeit.Email(),
		PhoneNumber:    gofakeit.Phone(),
		Password:       password,
		IsActive:       true,
	}
	createUser, err := authRepo.CreateUser(user)
	if err != nil {
		return nil, nil, err
	}

	accountReq := models.Account{
		BaseModel: models.BaseModel{
			ID:        uuid.New().String(),
			CreatedAt: accountTimeCreated,
			UpdatedAt: &accountTimeCreated,
		},
		UserID:           createUser.ID,
		AccountNumber:    extractAccountNumberFromPhoneNumber(createUser.PhoneNumber),
		AccountType:      models.AccountTypeSavings,
		Active:           true,
		TotalBalance:     0,
		AvailableBalance: 0,
		PendingBalance:   0,
		LockedBalance:    0,
	}
	account, err := walletRepo.CreateAccount(&accountReq)
	if err != nil {
		return nil, nil, err
	}

	return user, account, nil
}
