package models

type Transaction struct {
	BaseModel
	AccountID        string `json:"account_id"`
	Entry            string `json:"entry"`
	Purpose          string `json:"purpose"`
	Description      string `json:"description"`
	Remark           string `json:"remark"`
	Status           string `json:"status"`
	BeneficiaryName  string `json:"beneficiary_name"`
	TotalBalance     int64  `json:"total_balance"`
	AvailableBalance int64  `json:"available_balance"`
	PendingBalance   int64  `json:"pending_balance"`
}
