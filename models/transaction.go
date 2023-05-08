package models

type Transaction struct {
	BaseModel
	AccountID        string `json:"account_id"`
	Entry            string `json:"entry"`
	Change           int64  `json:"change"`
	Purpose          string `json:"purpose"`
	Status           string `json:"status"`
	TotalBalance     int64  `json:"total_balance"`
	AvailableBalance int64  `json:"available_balance"`
	PendingBalance   int64  `json:"pending_balance"`
	Reference        string `json:"reference"`
}
