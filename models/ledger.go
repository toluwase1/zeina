package models

type Ledger struct {
	ID           string `json:"id"`
	CreatedAt    int64  `json:"created_at"`
	AccountID    string `json:"account_id"` //fk
	CurrencyCode string `json:"currency_code"`
	Entry        string `json:"entry"`
	Change       int64  `json:"change"`
	Type         string `json:"type"` //:enum withdrawal, deposit eytc
}
