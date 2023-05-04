package models

type BlackList struct {
	BaseModel
	Token string `json:"token"`
	Email string `json:"email"`
}
