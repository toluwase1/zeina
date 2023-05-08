package services

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"regexp"
	"time"
	"zeina/config"
	"zeina/db"
	apiError "zeina/errors"
	"zeina/models"
	"zeina/services/jwt"
)

// AuthService interface
//
//go:generate mockgen -destination=../mocks/auth_mock.go -package=mocks zeina/services AuthService
type AuthService interface {
	LoginUser(request *models.LoginRequest) (*models.LoginResponse, *apiError.Error)
	SignupUser(request *models.User) (*models.User, *apiError.Error)
}

// authService struct
type authService struct {
	Config     *config.Config
	authRepo   db.AuthRepository
	walletRepo db.WalletRepository
}

// NewAuthService instantiate an authService
func NewAuthService(authRepo db.AuthRepository, walletRepo db.WalletRepository, conf *config.Config) AuthService {
	return &authService{
		Config:     conf,
		authRepo:   authRepo,
		walletRepo: walletRepo,
	}
}

func (a *authService) SignupUser(user *models.User) (*models.User, *apiError.Error) {
	err := a.authRepo.IsEmailExist(user.Email)
	if err != nil {
		return nil, apiError.New("email already exist", http.StatusBadRequest)
	}

	err = a.authRepo.IsPhoneExist(user.PhoneNumber)
	if err != nil {
		return nil, apiError.New("phone already exist", http.StatusBadRequest)
	}

	user.HashedPassword, err = GenerateHashPassword(user.Password)
	if err != nil {
		log.Printf("error generating password hash: %v", err.Error())
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}
	//CREATE USER
	user.ID = uuid.New().String()
	timeCreated := time.Now().Unix()
	user.CreatedAt = timeCreated
	user.UpdatedAt = &timeCreated
	user.Password = ""
	user.IsActive = true
	userCreated, err := a.authRepo.CreateUser(user)
	if err != nil {
		log.Printf("unable to create user: %v", err.Error())
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}

	//CREATE MULTIPLE USER ACCOUNTS FOR ALL ACCOUNT TYPES WITH 0 BALANCES EACH AT SIGNUP
	bankAccountNumber := extractAccountNumberFromPhoneNumber(userCreated.PhoneNumber)
	accountTimeCreated := time.Now().Unix()
	accountReq := models.Account{
		BaseModel: models.BaseModel{
			ID:        uuid.New().String(),
			CreatedAt: accountTimeCreated,
			UpdatedAt: &accountTimeCreated,
		},
		UserID:           userCreated.ID,
		AccountNumber:    bankAccountNumber,
		AccountType:      models.AccountTypeSavings,
		Active:           true,
		TotalBalance:     0,
		AvailableBalance: 0,
		PendingBalance:   0,
		LockedBalance:    0,
	}
	_, err = a.walletRepo.CreateAccount(&accountReq)
	if err != nil {
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}
	//lockAccount := models.LockedBalance{
	//	BaseModel: models.BaseModel{
	//		ID:        uuid.New().String(),
	//		CreatedAt: accountTimeCreated,
	//		UpdatedAt: &accountTimeCreated,
	//	},
	//	AccountID:    acct.ID,
	//	LockDate:     0,
	//	ReleaseDate:  0,
	//	AmountLocked: 0,
	//}
	//err = a.walletRepo.CreateLockedAccount(&lockAccount)
	//if err != nil {
	//	return nil, apiError.New("internal server error", http.StatusInternalServerError)
	//}
	return userCreated, nil
}

func GenerateHashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func (a *authService) LoginUser(loginRequest *models.LoginRequest) (*models.LoginResponse, *apiError.Error) {
	foundUser, err := a.authRepo.FindUserByEmail(loginRequest.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apiError.New("invalid email", http.StatusUnprocessableEntity)
		} else {
			log.Printf("error from database: %v", err)
			return nil, apiError.ErrInternalServerError
		}
	}

	if foundUser.IsActive == false {
		return nil, apiError.New("email not verified", http.StatusUnauthorized)
	}

	if err := foundUser.VerifyPassword(loginRequest.Password); err != nil {
		return nil, apiError.ErrInvalidPassword
	}

	accessToken, err := jwt.GenerateToken(foundUser.Email, a.Config.JWTSecret)
	if err != nil {
		log.Printf("error generating token %s", err)
		return nil, apiError.ErrInternalServerError
	}

	return foundUser.LoginUserToDto(accessToken), nil
}

func extractAccountNumberFromPhoneNumber(phoneNumber string) string {
	pattern := regexp.MustCompile(`\d{10}$`)
	accountNumber := pattern.FindString(phoneNumber)
	return accountNumber
}
