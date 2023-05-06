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
	InternalMove(ctx context.Context, ledger models.Ledger, transaction models.Transaction, zeinaAccount models.Account, depositRequest models.DepositRequest) error
	//createLedgerRecord(ledger *models.Ledger) (*models.Ledger, error)
	//CreateZeinaLedgerRecord(ledger *models.Ledger, account *models.Account) error
	//CreateTransactionRecord(transaction models.Transaction) error
	//CreateZeinaTransactionRecord(transaction models.Transaction, account *models.Account) error
	//UpdateCustomerAccountRecord(amount int64, accountNumber string) (*models.Account, error)
	//UpdateZeinaAccountRecord(amount int64, number string) (*models.Account, error)
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

func createLedgerRecord(ctx context.Context, ledger models.Ledger, tx *sql.Tx) error {
	query := "INSERT INTO ledgers (id, created_at, account_id, account_type, entry, change, type) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, account_id, account_type, entry, change, type"
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("internal server error %v", err)
	}
	defer stmt.Close()
	log.Println("CreateLedgerRecord")
	var createdLedger models.Ledger
	err = stmt.QueryRowContext(ctx, ledger.ID, ledger.CreatedAt, ledger.AccountID, ledger.AccountType, ledger.Entry, ledger.Change, ledger.Type).
		Scan(&createdLedger.ID, &createdLedger.CreatedAt, &createdLedger.AccountID, &createdLedger.AccountType, &createdLedger.Entry, &createdLedger.Change, &createdLedger.Type)
	if err != nil {
		return fmt.Errorf("could not create ledger record: %v", err)
	}
	return nil
}

func (a *accountRepo) createZeinaLedgerRecord(ctx context.Context, ledger models.Ledger, account models.Account, tx *sql.Tx) error {

	zeinaAccount, err := a.FindZeinaAccount(a.Config.ZeinaAccountNumber)
	if err != nil {
		return err
	}
	ledger.ID = uuid.New().String()
	ledger.AccountID = zeinaAccount.ID
	query := "INSERT INTO ledgers (id, created_at, account_id, account_type, entry, change, type) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, account_id, account_type, entry, change, type"
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("internal server error %v", err)
	}
	_, err = tx.ExecContext(ctx, query, ledger.ID, ledger.AccountID, ledger.Entry, ledger.Purpose, transaction.Description, transaction.Remark, transaction.Status, transaction.BeneficiaryName, transaction.TotalBalance, transaction.AvailableBalance, transaction.PendingBalance, transaction.CreatedAt, transaction.UpdatedAt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return nil
}

func (a *accountRepo) createTransactionRecord(ctx context.Context, transaction models.Transaction, tx *sql.Tx) error {
	log.Println("CreateTransactionRecord")
	query := `INSERT INTO transactions (id, account_id, entry, purpose, description, remark, status, beneficiary_name, total_balance, available_balance, pending_balance, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := tx.ExecContext(ctx, query, transaction.ID, transaction.AccountID, transaction.Entry, transaction.Purpose, transaction.Description, transaction.Remark, transaction.Status, transaction.BeneficiaryName, transaction.TotalBalance, transaction.AvailableBalance, transaction.PendingBalance, transaction.CreatedAt, transaction.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (a *accountRepo) CreateZeinaTransactionRecord(ctx context.Context, transaction models.Transaction, account *models.Account, tx *sql.Tx) error {
	log.Println("CreateZeinaTransactionRecord")
	zeinaAccount, err := a.FindZeinaAccount(a.Config.ZeinaAccountNumber)
	if err != nil {
		return err
	}
	transaction.AccountID = zeinaAccount.ID
	query := `INSERT INTO transactions (id, account_id, entry, purpose, description, remark, status, beneficiary_name, total_balance, available_balance, pending_balance, created_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	tx, err = a.DB.Begin()
	if err != nil {
		return err
	}

	// Lock the account for the transaction
	lockQuery := `SELECT * FROM accounts WHERE account_number = $1 FOR UPDATE`
	row := tx.QueryRowContext(ctx, lockQuery, account.AccountNumber)
	err = row.Scan(&account.ID, &account.UserID, &account.AccountNumber, &account.AccountType, &account.Active, &account.TotalBalance, &account.AvailableBalance, &account.PendingBalance, &account.LockedBalance, &account.CreatedAt, &account.UpdatedAt, &account.DeletedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update the account balance
	updateQuery := `UPDATE accounts SET total_balance = total_balance + $1, available_balance = available_balance + $1, locked_balance = locked_balance - $1 WHERE account_number = $2`
	_, err = tx.Exec(updateQuery, transaction.TotalBalance, account.AccountNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create the transaction record
	_, err = tx.Exec(query, transaction.ID, transaction.AccountID, transaction.Entry, transaction.Purpose, transaction.Description, transaction.Remark, transaction.Status, transaction.BeneficiaryName, transaction.TotalBalance, transaction.AvailableBalance, transaction.PendingBalance, transaction.CreatedAt)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (a *accountRepo) updateCustomerAccountRecord(ctx context.Context, amount int64, accountNumber string, tx *sql.Tx) error {
	number, err := a.FindAccountByNumber(accountNumber)
	if err != nil {
		return err
	}
	log.Println("check", number)
	log.Println("UpdateCustomerAccountRecord", accountNumber)
	query := "UPDATE accounts SET pending_balance = pending_balance + $1 WHERE account_number = $2"

	res, err := tx.ExecContext(ctx, query, amount, accountNumber)
	if err != nil {
		return err
	}
	row, _ := res.RowsAffected()
	if row < 1 {
		return fmt.Errorf("could not update account")
	}

	return nil
}

func (a *accountRepo) updateZeinaAccountRecord(ctx context.Context, amount int64, number string, tx *sql.Tx) error {
	log.Println("UpdateZeinaAccountRecord")
	zeinaAccount, err := a.FindZeinaAccount(a.Config.ZeinaAccountNumber)
	if err != nil {
		return err
	}
	query := "UPDATE accounts SET total_balance = total_balance - $1, available_balance = available_balance - $1, locked_balance = locked_balance + $1 WHERE id = $2"
	res, err := tx.ExecContext(ctx, query, amount, zeinaAccount.ID)
	if err != nil {
		return err
	}
	row, _ := res.RowsAffected()
	if row < 1 {
		return fmt.Errorf("could not update account")
	}

	return nil
}

//err = tx.QueryRow("UPDATE my_table SET my_column = $1 WHERE id = $2 RETURNING my_column", "new value", 123).Scan(&updatedValue)

func (a *accountRepo) InternalMove(ctx context.Context, ledger models.Ledger, transaction models.Transaction, zeinaAccount models.Account, depositRequest models.DepositRequest) error {
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}

	errorGroup, ctx := errgroup.WithContext(ctx)
	errorGroup.Go(func() error {
		return createLedgerRecord(ctx, ledger, tx)
	})

	errorGroup.Go(func() error {
		return a.createZeinaLedgerRecord(ctx, ledger, zeinaAccount, tx)

	})
	errorGroup.Go(func() error {
		return a.createTransactionRecord(ctx, transaction, tx)
	})
	errorGroup.Go(func() error {
		return a.updateCustomerAccountRecord(ctx, depositRequest.Amount, depositRequest.AccountNumber, tx)

	})
	errorGroup.Go(func() error {
		return a.updateZeinaAccountRecord(ctx, depositRequest.Amount, depositRequest.AccountNumber, tx)

	})
	//if err = errorGroup.Wait(); err != nil {
	//	rollBackErr := tx.Rollback()
	//	if rollBackErr != nil {
	//		return rollBackErr
	//	}
	//	return err
	//}
	//txErr := tx.Commit()
	//if err != nil {
	//	return txErr
	//}
	return nil
}
