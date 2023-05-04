package models

type User struct {
	BaseModel
	Email string `json:"email"`
	Name  string `json:"name"`
}
