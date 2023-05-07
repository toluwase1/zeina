package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"log"
	"zeina/config"
	"zeina/models"
)

//go:generate mockgen -destination=../mocks/wallet_repo_mock.go -package=mocks zeina/db WalletRepository
type WalletRepository interface {
	CreateAccount(account *models.Account) (*models.Account, error)
	FindAccountByNumber(id string) (*models.Account, error)
	FindAccountByUserID(id string) (*models.Account, error)
	FindZeinaAccount(id string) (models.Account, error)
	FindAccountByType(accountType string) (*models.Account, error)
	InternalMove(ctx context.Context, ledger models.Ledger, transaction models.Transaction) error //createLedgerRecord(ledger *models.Ledger) (*models.Ledger, error)
	ExternalMove(ctx context.Context, ledger models.Ledger, transaction models.Transaction) error
}

type accountRepo struct {
	DB     *sql.DB
	Config *config.Config
}

func NewWalletRepo(db *SqlDB, conf *config.Config) WalletRepository {
	return &accountRepo{DB: db.DB, Config: conf}
}

func (a *accountRepo) CreateAccount(account *models.Account) (*models.Account, error) {
	query := "INSERT INTO accounts (id, user_id, active, account_number, account_type, total_balance, available_balance, pending_balance, locked_balance, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10 ,$11) RETURNING id, user_id, active, account_number, account_type, total_balance, available_balance, pending_balance, locked_balance, created_at, updated_at"
	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("internal server error %v", err)
	}
	defer stmt.Close()

	var createdAccount models.Account
	err = stmt.QueryRow(account.ID, account.UserID, account.Active, account.AccountNumber, account.AccountType, account.TotalBalance, account.AvailableBalance, account.PendingBalance, account.LockedBalance, account.CreatedAt, account.UpdatedAt).
		Scan(&createdAccount.ID, &createdAccount.UserID, &createdAccount.Active, &createdAccount.AccountNumber, &createdAccount.AccountType, &createdAccount.TotalBalance, &createdAccount.AvailableBalance, &createdAccount.PendingBalance, &createdAccount.LockedBalance, &createdAccount.CreatedAt, &createdAccount.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not create account: %v", err)
	}
	return &createdAccount, nil
}

func (a *accountRepo) FindAccountByNumber(accountNumber string) (*models.Account, error) {
	query := "SELECT id, user_id, active, account_number, account_type, total_balance, available_balance, pending_balance, locked_balance, created_at, updated_at FROM accounts WHERE account_number=$1"
	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("internal server error: %v", err)
	}
	defer stmt.Close()

	var account models.Account
	err = stmt.QueryRow(accountNumber).Scan(&account.ID, &account.UserID, &account.Active, &account.AccountNumber, &account.AccountType, &account.TotalBalance, &account.AvailableBalance, &account.PendingBalance, &account.LockedBalance, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	return &account, nil
}

func (a *accountRepo) FindZeinaAccount(accountNumber string) (models.Account, error) {
	query := "SELECT id, user_id, active, account_number, account_type, total_balance, available_balance, pending_balance, locked_balance, created_at, updated_at FROM accounts WHERE account_number=$1"
	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return models.Account{}, fmt.Errorf("internal server error: %v", err)
	}
	defer stmt.Close()

	var account models.Account
	err = stmt.QueryRow(accountNumber).Scan(&account.ID, &account.UserID, &account.Active, &account.AccountNumber, &account.AccountType, &account.TotalBalance, &account.AvailableBalance, &account.PendingBalance, &account.LockedBalance, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Account{}, fmt.Errorf("account not found")
		}
		return models.Account{}, fmt.Errorf("internal server error: %v", err)
	}

	return account, nil
}

func (a *accountRepo) FindAccountByUserID(id string) (*models.Account, error) {
	query := "SELECT id, user_id, active, account_number, account_type, total_balance, available_balance, pending_balance, locked_balance, created_at, updated_at FROM accounts WHERE user_id=$1"
	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("internal server error: %v", err)
	}
	defer stmt.Close()

	var account models.Account
	err = stmt.QueryRow(id).Scan(&account.ID, &account.UserID, &account.Active, &account.AccountNumber, &account.AccountType, &account.TotalBalance, &account.AvailableBalance, &account.PendingBalance, &account.LockedBalance, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	return &account, nil
}

func (a *accountRepo) FindAccountByType(accountType string) (*models.Account, error) {
	query := "SELECT id, user_id, active, account_number, account_type, total_balance, available_balance, pending_balance, locked_balance, created_at, updated_at FROM accounts WHERE account_type=$1"
	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("internal server error: %v", err)
	}
	defer stmt.Close()

	var account models.Account
	err = stmt.QueryRow(accountType).Scan(&account.ID, &account.UserID, &account.Active, &account.AccountNumber, &account.AccountType, &account.TotalBalance, &account.AvailableBalance, &account.PendingBalance, &account.LockedBalance, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	return &account, nil
}

func (a *accountRepo) createLedgerRecord(ctx context.Context, ledger models.Ledger, tx *sql.Tx) error {
	query := "INSERT INTO ledgers (id, created_at, account_id, account_type, entry, change, type) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, account_id, account_type, entry, change, type"
	log.Println("CreateLedgerRecord")
	for _, entry := range ledger.Entries {
		_, err := tx.ExecContext(ctx, query, uuid.New().String(), ledger.CreatedAt, entry.AccountID, ledger.AccountType, entry.Side, entry.Delta, ledger.Type)
		if err != nil {
			return fmt.Errorf("could not create ledger record: %v", err)
		}
	}
	return nil
}
func (a *accountRepo) createTransactionRecord(ctx context.Context, transaction models.Transaction, tx *sql.Tx) error {
	log.Println("CreateTransactionRecord")
	query := `INSERT INTO transactions (id, account_id, entry, purpose, status, available_balance, pending_balance, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := tx.ExecContext(ctx, query, transaction.ID, transaction.AccountID, transaction.Entry, transaction.Purpose, transaction.Status, transaction.AvailableBalance, transaction.PendingBalance, transaction.CreatedAt, transaction.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (a *accountRepo) externalCustomerAccountRecord(ctx context.Context, amount int64, accountNumber string, tx *sql.Tx) error {
	query := "UPDATE accounts SET pending_balance = pending_balance - $1, available_balance = available_balance + $2 WHERE account_number = $2"

	res, err := tx.ExecContext(ctx, query, amount, accountNumber)
	if err != nil {
		return err
	}
	row, _ := res.RowsAffected()
	if row < 1 {
		return fmt.Errorf("could not update account kkkkk")
	}

	return nil
}

func (a *accountRepo) updateCustomerAccountRecord(ctx context.Context, amount int64, accountId string, tx *sql.Tx) error {
	query := "UPDATE accounts SET pending_balance = pending_balance + $1 WHERE id = $2"

	res, err := tx.ExecContext(ctx, query, amount, accountId)
	if err != nil {
		return err
	}
	row, _ := res.RowsAffected()
	if row < 1 {
		return fmt.Errorf("not update account")
	}

	return nil
}

func (a *accountRepo) InternalMove(ctx context.Context, ledger models.Ledger, transaction models.Transaction) error {
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}

	errorGroup, ctx := errgroup.WithContext(ctx)
	errorGroup.Go(func() error {
		return a.createLedgerRecord(ctx, ledger, tx)
	})
	errorGroup.Go(func() error {
		return a.updateCustomerAccountRecord(ctx, ledger.Entries[0].Delta, ledger.Entries[0].AccountID, tx)
	})
	errorGroup.Go(func() error {
		return a.createTransactionRecord(ctx, transaction, tx)
	})
	if err = errorGroup.Wait(); err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return rollBackErr
		}
		return err
	}
	txErr := tx.Commit()
	if err != nil {
		return txErr
	}
	return nil
}

func (a *accountRepo) ExternalMove(ctx context.Context, ledger models.Ledger, transaction models.Transaction) error {
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}

	errorGroup, ctx := errgroup.WithContext(ctx)
	errorGroup.Go(func() error {
		return a.createLedgerRecord(ctx, ledger, tx)
	})
	errorGroup.Go(func() error {
		return a.externalCustomerAccountRecord(ctx, ledger.Entries[0].Delta, ledger.Entries[0].AccountID, tx)
	})
	errorGroup.Go(func() error {
		return a.createTransactionRecord(ctx, transaction, tx)
	})
	if err = errorGroup.Wait(); err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return rollBackErr
		}
		return err
	}
	txErr := tx.Commit()
	if err != nil {
		return txErr
	}
	return nil
}
