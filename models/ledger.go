package models

const (
	Deposit       = "DEPOSIT"
	Withdrawal    = "WITHDRAWAL"
	StatusPending = "pending"
	StatusSuccess = "success"
)

type Ledger struct {
	ID          string `json:"id"`
	CreatedAt   int64  `json:"created_at"`
	AccountID   string `json:"account_id"`
	AccountType string `json:"account_type"`
	Entry       string `json:"entry"`
	Change      int64  `json:"change"`
	Type        string `json:"type"`
}
