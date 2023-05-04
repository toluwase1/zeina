package db

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"zeina/models"
)

// DB provides access to the different db
//go:generate mockgen -destination=../mocks/auth_repo_mock.go -package=mocks github.com/decagonhq/zeina/db AuthRepository
convert all my functions below functions to raw sql queries from gorm queries, ensure that everything is correct:
type AuthRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	IsEmailExist(email string) error
	IsPhoneExist(email string) error
	FindUserByUsername(username string) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	AddToBlackList(blacklist *models.BlackList) error
	TokenInBlacklist(token string) bool
	IsTokenInBlacklist(token string) error
	UpdatePassword(password string, email string) error
}

type authRepo struct {
	DB *sql.DB
}

func NewAuthRepo(db *SqlDB) AuthRepository {
	return &authRepo{db.DB}
}
func (a *authRepo) CreateUser(user *models.User) (*models.User, error) {
	query := "INSERT INTO users (name, email, phone_number, hashed_password, is_email_active) VALUES ($1, $2, $3, $4, $5) RETURNING id, username, email, phone_number, hashed_password, is_email_active, created_at, updated_at"
	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %v", err)
	}
	defer stmt.Close()

	var createdUser models.User
	err = stmt.QueryRow(user.Name, user.Email, user.PhoneNumber, user.HashedPassword, user.IsEmailActive).Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email, &createdUser.PhoneNumber, &createdUser.HashedPassword, &createdUser.IsEmailActive, &createdUser.CreatedAt, &createdUser.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not create user: %v", err)
	}
	return &createdUser, nil
}

func (a *authRepo) FindUserByUsername(name string) (*models.User, error) {
	query := "SELECT * FROM users WHERE email = $1 OR name = $1"
	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %v", err)
	}
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRow(name).Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.HashedPassword, &user.IsEmailActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}
	return &user, nil
}

func (a *authRepo) IsEmailExist(email string) error {
	query := "SELECT COUNT(*) FROM users WHERE email = $1"
	var count int64
	err := a.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return fmt.Errorf("error in counting: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("email already in use")
	}
	return nil
}

func (a *authRepo) IsPhoneExist(phone string) error {
	query := "SELECT COUNT(*) FROM users WHERE phone_number = $1"
	var count int64
	err := a.DB.QueryRow(query, phone).Scan(&count)
	if err != nil {
		return fmt.Errorf("error in counting: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("phone number already in use")
	}
	return nil
}

func (a *authRepo) FindUserByEmail(email string) (*models.User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %v", err)
	}
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRow(email).Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.HashedPassword, &user.IsEmailActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *authRepo) UpdateUser(user *models.User) error {
	query := "UPDATE users SET name = ?, email = ?, phone_number = ? WHERE id = ?"
	result, err := a.DB.Exec(query, user.Name, user.Email, user.PhoneNumber, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (a *authRepo) TokenInBlacklist(token string) bool {
	var count int64
	query := "SELECT COUNT(*) FROM blacklists WHERE token = ?"
	err := a.DB.QueryRow(query, token).Scan(&count)
	if err != nil {
		return true
	}
	return count > 0
}

func (a *authRepo) IsTokenInBlacklist(token string) error {
	var count int64
	query := "SELECT COUNT(*) FROM blacklists WHERE token = ?"
	err := a.DB.QueryRow(query, token).Scan(&count)
	if err != nil {
		return errors.Wrap(err, "sql query error")
	}
	if count > 0 {
		return fmt.Errorf("token expired, request a new link")
	}
	return nil
}


func (a *authRepo) AddToBlackList(blacklist *models.BlackList) error {
	_, err := a.DB.Exec("INSERT INTO blacklist (email, token) VALUES (?, ?)", blacklist.Email, blacklist.Token)
	if err != nil {
		return err
	}
	return nil
}


func (a *authRepo) UpdatePassword(password string, email string) error {
	_, err := a.DB.Exec("UPDATE users SET password = ? WHERE email = ?", password, email)
	if err != nil {
		return err
	}
	return nil
}
