package models

import "time"

type TransferDirection string
type TransferStatus string
type BaseModel struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt *int64 `json:"updated_at"`
	DeletedAt *int64 `json:"deleted_at"`
}

type Webhook struct {
	Transaction string      `json:"transaction"`
	Time        time.Time   `json:"time"`
	Data        interface{} `json:"data"`
}

type OutgoingWebhookPayload struct {
	Notify     string       `json:"notify"`
	NotifyType string       `json:"notifyType"`
	Data       OutgoingData `json:"data"`
}

type OutgoingData struct {
	Id            int    `json:"id"`
	Reference     string `json:"reference"`
	Sessionid     string `json:"sessionid"`
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	Fee           string `json:"fee"`
	BankCode      string `json:"bank_code"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	CountryCode   string `json:"countryCode"`
	PaymentMode   string `json:"paymentMode"`
	Narration     string `json:"narration"`
	Sender        string `json:"sender"`
	Domain        string `json:"domain"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type LedgerTransaction struct {
	LedgerID        string            `json:"ledger_id,omitempty"`
	Ledger          Ledger            `json:"ledger,omitempty" gorm:"foreignKey:LedgerID"`
	Direction       TransferDirection `json:"direction"`
	Amount          int64             `json:"amount"`
	Fee             int64             `json:"fee"`
	ServiceFee      int64             `json:"-"`
	CustomerFee     int64             `json:"-"`
	CustomerName    string            `json:"customer_name"`
	CurrencyCode    string            `json:"currency_code"`
	Status          TransferStatus    `json:"status,omitempty" gorm:"default:'pending'"`
	BeneficiaryName string            `json:"beneficiary_name"`
	TransferService string            `json:"-"`
	Reference       string            `json:"reference"`
	BankName        string            `json:"bank_name"`
	BankCode        string            `json:"bank_code"`
	Narration       string            `json:"narration"`
	TransferType    string            `json:"transfer_type"`
}
