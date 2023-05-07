package models

const (
	Deposit         = "DEPOSIT"
	Withdrawal      = "WITHDRAWAL"
	StatusPending   = "pending"
	StatusCompleted = "completed"
)

type Ledger struct {
	ID          string  `json:"id"`
	CreatedAt   int64   `json:"created_at"`
	AccountType string  `json:"account_type"`
	Entries     []Entry `json:"entry"`
	Type        string  `json:"type"`
}

type Entry struct {
	AccName   string `json:"acc_name"`
	AccountID string `json:"account_id"`
	Delta     int64  `json:"delta"`
	Side      string `json:"side"`
}
