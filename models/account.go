package models

const (
	AccountTypeSavings = "SAVINGS"
	CreditEntry        = "credit"
	DebitEntry         = "debit"
)

type Account struct {
	BaseModel
	UserID           string `json:"user_id"`
	AccountNumber    string `json:"account_number"`
	AccountType      string `json:"account_type"`
	Active           bool   `json:"active"`
	TotalBalance     int64  `json:"total_balance"`
	AvailableBalance int64  `json:"available_balance"`
	PendingBalance   int64  `json:"pending_balance"`
	LockedBalance    int64  `json:"locked_balance"`
}

type LockedBalance struct {
	BaseModel
	AccountID    string `json:"account_id"`
	LockDate     int64  `json:"lock_date"`
	ReleaseDate  int64  `json:"release_date"`
	AmountLocked int64  `json:"amount_locked"`
}

type DepositRequest struct {
	Amount        int64  `json:"amount"`
	AccountType   string `json:"account_type"`
	AccountNumber string `json:"account_number"`
}
