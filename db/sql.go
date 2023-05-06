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
	seedZeinasAccount(sql.DB)
}

func getPostgresDB(c *config.Config) *sql.DB {
	log.Printf("Connecting to postgres: %+v", c)
	postgresDSN := "postgres://postgres:toluwase@localhost/zeina?sslmode=disable"
	//postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d TimeZone=Africa/Lagos",
	//	c.PostgresHost, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresPort)
	log.Println(postgresDSN)
	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		log.Println("db connection error", err)
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Connected to database successfully!")
	return db
}

func seedZeinasAccount(DB *sql.DB) {
	user := models.User{}
	user.ID = uuid.New().String()
	user.Name = "zeina"
	user.PhoneNumber = "+23481111111111"
	user.Email = "zeina@gmail.com"
	timeCreated := time.Now().Unix()
	user.CreatedAt = timeCreated
	user.UpdatedAt = &timeCreated
	user.Password = ""
	user.IsActive = true
	accountTimeCreated := time.Now().Unix()
	accountReq := models.Account{
		BaseModel: models.BaseModel{
			ID:        uuid.New().String(),
			CreatedAt: accountTimeCreated,
			UpdatedAt: &accountTimeCreated,
		},
		UserID:           user.ID,
		AccountNumber:    "1111111111",
		AccountType:      models.AccountTypeSavings,
		Active:           true,
		TotalBalance:     0,
		AvailableBalance: 0,
		PendingBalance:   0,
		LockedBalance:    0,
	}
	err := seedUser(&user, DB)
	if err != nil {
		return
	}
	err = createAccount(&accountReq, DB)
	if err != nil {
		return
	}
}

func createAccount(account *models.Account, DB *sql.DB) error {
	query := "INSERT INTO accounts (id, user_id, active, account_number, account_type, total_balance, available_balance, pending_balance, locked_balance, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10 ,$11) RETURNING id, user_id, active, account_number, account_type, total_balance, available_balance, pending_balance, locked_balance, created_at, updated_at"
	stmt, err := DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("internal server error %v", err)
	}
	defer stmt.Close()

	var createdAccount models.Account
	err = stmt.QueryRow(account.ID, account.UserID, account.Active, account.AccountNumber, account.AccountType, account.TotalBalance, account.AvailableBalance, account.PendingBalance, account.LockedBalance, account.CreatedAt, account.UpdatedAt).
		Scan(&createdAccount.ID, &createdAccount.UserID, &createdAccount.Active, &createdAccount.AccountNumber, &createdAccount.AccountType, &createdAccount.TotalBalance, &createdAccount.AvailableBalance, &createdAccount.PendingBalance, &createdAccount.LockedBalance, &createdAccount.CreatedAt, &createdAccount.UpdatedAt)
	if err != nil {
		return fmt.Errorf("could not create account: %v", err)
	}
	log.Println("successful")
	return nil
}

func seedUser(user *models.User, DB *sql.DB) error {
	query := "INSERT INTO users (id, name, email, phone_number, hashed_password, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, name, email, phone_number, hashed_password, is_active, created_at, updated_at"
	stmt, err := DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("internal server error %v", err)
	}
	defer stmt.Close()

	var createdUser models.User
	err = stmt.QueryRow(user.ID, user.Name, user.Email, user.PhoneNumber, user.HashedPassword, user.IsActive, user.CreatedAt, user.UpdatedAt).Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email, &createdUser.PhoneNumber, &createdUser.HashedPassword, &createdUser.IsActive, &createdUser.CreatedAt, &createdUser.UpdatedAt)
	if err != nil {
		return fmt.Errorf("could not create user: %v", err)
	}
	return nil
}
