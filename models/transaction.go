package models

type Transaction struct {
	Model
	AccountId      string `json:"account_id"`      //fk
	AccountBalance int64  `json:"account_balance"` //currentaccountbal noq
	//cr/dbjson:Side[cr,db]
	Entry           string `json:"entry"`
	Purpose         string `json:"purpose"`
	Description     string `json:"description"`
	Remark          string `json:"remark"`
	Status          string `json:"status"`
	BeneficiaryName string `json:"beneficiary_name"`
	//customerName

}
