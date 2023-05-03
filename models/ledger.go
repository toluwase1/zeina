package models

type Ledger struct {
	ID           string `json:"id"`
	CreatedAt    int64  `json:"created_at"`
	AccountId    string `json:"account_id"` //fk
	CurrencyCode string `json:"currency_code"`
	Entry        string `json:"entry"`
	//change/beta/amount
	Change int64  `json:"change"`
	Reason string `json:"reason"` //:enum
	//Noupdatedat at all

}
