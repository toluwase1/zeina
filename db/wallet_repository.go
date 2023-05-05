package db

import (
	"database/sql"
	"fmt"
	"zeina/models"
)

//go:generate mockgen -destination=../mocks/wallet_repo_mock.go -package=mocks zeina/db WalletRepository
type WalletRepository interface {
	CreateAccount(account *models.Account) (*models.Account, error)
}

type accountRepo struct {
	DB *sql.DB
}

func NewWalletRepo(db *SqlDB) WalletRepository {
	return &accountRepo{db.DB}
}

func (a *accountRepo) CreateAccount(account *models.Account) (*models.Account, error) {
	query := "INSERT INTO accounts (id, user_id, active, account_number, total_balance, available_balance, pending_balance, locked_balance, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, user_id, active, account_number, total_balance, available_balance, pending_balance, locked_balance, created_at, updated_at"
	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("internal server error %v", err)
	}
	defer stmt.Close()

	var createdAccount models.Account
	err = stmt.QueryRow(account.ID, account.UserID, account.Active, account.AccountNumber, account.TotalBalance, account.AvailableBalance, account.PendingBalance, account.LockedBalance, account.CreatedAt, account.UpdatedAt).
		Scan(&createdAccount.ID, &createdAccount.UserID, &createdAccount.Active, &createdAccount.AccountNumber, &createdAccount.TotalBalance, &createdAccount.AvailableBalance, &createdAccount.PendingBalance, &createdAccount.LockedBalance, &createdAccount.CreatedAt, &createdAccount.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not create account: %v", err)
	}
	return &createdAccount, nil
}
