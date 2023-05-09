package db

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"time"
	"zeina/config"
	"zeina/models"
)

type SqlDB struct {
	DB *sql.DB
}

func GetDB(c *config.Config) *SqlDB {
	sqlDB := &SqlDB{}
	sqlDB.Init(c)
	return sqlDB
}

func (sql *SqlDB) Init(c *config.Config) {
	sql.DB = getPostgresDB(c)
	err := createTables(sql.DB)
	if err != nil {
		log.Println("error creating tables: ", err)
		return
	} else {
		log.Println("created tables successfully: ", err)
	}
	err = seedZeinasAccount(sql.DB)
	if err != nil {
		log.Println("seeder error: ", err)
		return
	}
}

func getPostgresDB(c *config.Config) *sql.DB {
	log.Printf("Connecting to postgres: %+v", c)
	postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d TimeZone=Africa/Lagos sslmode=disable",
		c.PostgresHost, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresPort)
	log.Println(postgresDSN)
	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		log.Println("db connection error", err)
		panic(err)
	}

	log.Println("Connected to database successfully!")
	return db
}

func createTables(DB *sql.DB) error {
	// create users table
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id uuid NOT NULL,
		email varchar(255) NOT NULL,
		name varchar(255) NOT NULL,
		phone_number varchar(255) NOT NULL,
		hashed_password varchar(255) NOT NULL,
		is_active varchar(255) NOT NULL,
		created_at bigint NOT NULL,
		updated_at bigint NOT NULL,
		deleted_at bigint DEFAULT NULL,
		PRIMARY KEY (id)
	)`)
	if err != nil {
		return err
	}

	// create accounts table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS accounts (
		id uuid NOT NULL,
		created_at bigint NOT NULL,
		updated_at bigint NOT NULL,
		deleted_at bigint DEFAULT NULL,
		user_id uuid NOT NULL,
		account_number varchar(255) NOT NULL,
		account_type varchar(255) NOT NULL,
		active boolean NOT NULL,
		total_balance bigint NOT NULL,
		available_balance bigint NOT NULL,
		pending_balance bigint NOT NULL,
		locked_balance bigint NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`)
	if err != nil {
		return err
	}

	// create ledgers table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS ledgers (
		id uuid NOT NULL,
		created_at bigint NOT NULL,
		account_id uuid NOT NULL,
		account_type varchar(255) NOT NULL,
		entry varchar(255) NOT NULL,
		change bigint NOT NULL,
		type varchar(255) NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (account_id) REFERENCES accounts(id)
	)`)
	if err != nil {
		return err
	}

	// create locked_balances table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS locked_balances (
		id uuid NOT NULL,
		created_at bigint NOT NULL,
		updated_at bigint NOT NULL,
		deleted_at bigint DEFAULT NULL,
		account_id uuid NOT NULL,
		lock_date bigint NOT NULL,
		release_date bigint NOT NULL,
		amount_locked bigint NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (account_id) REFERENCES accounts(id)
	)`)
	if err != nil {
		return err
	}

	// create transactions table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS transactions (
		id uuid NOT NULL,
		created_at bigint NOT NULL,
		updated_at bigint NOT NULL,
		deleted_at bigint DEFAULT NULL,
		account_id uuid NOT NULL,
		entry varchar(255) NOT NULL,
		purpose varchar(255) NOT NULL,
		status varchar(255) NOT NULL,
		change bigint DEFAULT NULL,
		available_balance bigint NOT NULL,
		pending_balance bigint NOT NULL,
		reference varchar(255) NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (account_id) REFERENCES accounts(id),
		CONSTRAINT unique_reference UNIQUE (reference)
	)`)
	if err != nil {
		return err
	}

	// create black_lists table
	_, err = DB.Exec(`CREATE TABLE black_lists (
                             id uuid NOT NULL,
                             created_at bigint NOT NULL,
                             updated_at bigint NOT NULL,
                             deleted_at bigint DEFAULT NULL,
                             token varchar(255) NOT NULL,
                             email varchar(255) NOT NULL,
                             PRIMARY KEY (id)
)`)
	if err != nil {
		return nil
	}

	return nil
}
func seedZeinasAccount(DB *sql.DB) error {
	user := models.User{}
	userID := uuid.New().String()
	user.ID = userID
	log.Println(userID)
	user.Name = "zeina"
	user.PhoneNumber = "+23481111111111"
	user.Email = "zeina@gmail.com"
	timeCreated := time.Now().Unix()
	user.CreatedAt = timeCreated
	user.UpdatedAt = &timeCreated
	user.Password = ""
	user.IsActive = true
	accountTimeCreated := time.Now().Unix()
	userr, err := seedAdminUser(&user, DB)
	if err != nil {
		log.Println("check if error is not nil: ", err)
		return err
	}
	log.Println("check if user was created: ", userr)
	accountReq := models.Account{
		BaseModel: models.BaseModel{
			ID:        uuid.New().String(),
			CreatedAt: accountTimeCreated,
			UpdatedAt: &accountTimeCreated,
		},
		UserID:           userr.ID,
		AccountNumber:    "1111111111",
		AccountType:      models.AccountTypeSavings,
		Active:           true,
		TotalBalance:     0,
		AvailableBalance: 0,
		PendingBalance:   0,
		LockedBalance:    0,
	}

	err = seedAdminAccount(&accountReq, DB)
	if err != nil {
		log.Println("error seeding", err)
		return err
	}
	return nil
}

func seedAdminAccount(account *models.Account, DB *sql.DB) error {
	query := "SELECT id FROM accounts WHERE account_number=$1"
	var id string
	row := DB.QueryRow(query, "1111111111")
	err := row.Scan(&id)
	if err == nil {
		// account already exists, return nil to indicate success
		return nil
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("error checking if account exists: %v", err)
	}

	query = "INSERT INTO accounts (id, user_id, active, account_number, account_type, total_balance, available_balance, pending_balance, locked_balance, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10 ,$11) RETURNING id, user_id, active, account_number, account_type, total_balance, available_balance, pending_balance, locked_balance, created_at, updated_at"
	stmt, err := DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("internal server error %v", err)
	}
	defer stmt.Close()

	var createdAccount models.Account
	err = stmt.QueryRow(account.ID, account.UserID, account.Active, account.AccountNumber, account.AccountType, account.TotalBalance, account.AvailableBalance, account.PendingBalance, account.LockedBalance, account.CreatedAt, account.UpdatedAt).
		Scan(&createdAccount.ID, &id, &createdAccount.Active, &createdAccount.AccountNumber, &createdAccount.AccountType, &createdAccount.TotalBalance, &createdAccount.AvailableBalance, &createdAccount.PendingBalance, &createdAccount.LockedBalance, &createdAccount.CreatedAt, &createdAccount.UpdatedAt)
	if err != nil {
		return fmt.Errorf("could not create account: %v", err)
	}
	log.Println("successful")
	return nil
}

//func seedAdminUser(user *models.User, DB *sql.DB) (*models.User, error) {
//	query := "SELECT id FROM users WHERE email=$1"
//	var id string
//	row := DB.QueryRow(query, user.Email)
//	err := row.Scan(&id)
//	if err == nil {
//		// user already exists, return nil to indicate success
//		return nil, err
//	} else if err != sql.ErrNoRows {
//		return nil, fmt.Errorf("error checking if user exists: %v", err)
//	}
//
//	query = "INSERT INTO users (id, name, email, phone_number, hashed_password, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, name, email, phone_number, hashed_password, is_active, created_at, updated_at"
//	stmt, err := DB.Prepare(query)
//	if err != nil {
//		return nil, fmt.Errorf("internal server error %v", err)
//	}
//	defer stmt.Close()
//
//	var createdUser models.User
//	err = stmt.QueryRow(user.ID, user.Name, user.Email, user.PhoneNumber, user.HashedPassword, user.IsActive, user.CreatedAt, user.UpdatedAt).Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email, &createdUser.PhoneNumber, &createdUser.HashedPassword, &createdUser.IsActive, &createdUser.CreatedAt, &createdUser.UpdatedAt)
//	if err != nil {
//		return nil, fmt.Errorf("could not create user: %v", err)
//	}
//	log.Println(createdUser)
//	return &createdUser, nil
//}

func seedAdminUser(user *models.User, DB *sql.DB) (*models.User, error) {

	query := "SELECT id FROM users WHERE email=$1"
	var id string
	row := DB.QueryRow(query, user.Email)
	err := row.Scan(&id)
	if err == nil {
		// user already exists, return the existing user
		existingUser := &models.User{
			BaseModel: models.BaseModel{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
			Name:           user.Name,
			Email:          user.Email,
			PhoneNumber:    user.PhoneNumber,
			HashedPassword: user.HashedPassword,
			IsActive:       user.IsActive,
		}
		return existingUser, nil
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("error checking if user exists: %v", err)
	}

	query = "INSERT INTO users (id, name, email, phone_number, hashed_password, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, name, email, phone_number, hashed_password, is_active, created_at, updated_at"
	stmt, err := DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("internal server error %v", err)
	}
	defer stmt.Close()

	var createdUser models.User
	err = stmt.QueryRow(user.ID, user.Name, user.Email, user.PhoneNumber, user.HashedPassword, user.IsActive, user.CreatedAt, user.UpdatedAt).Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email, &createdUser.PhoneNumber, &createdUser.HashedPassword, &createdUser.IsActive, &createdUser.CreatedAt, &createdUser.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not create user: %v", err)
	}
	log.Println("whats up", createdUser.ID, user.ID)
	return &createdUser, nil
}
