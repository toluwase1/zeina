package models

type User struct {
	BaseModel
	Email          string `json:"email"`
	Name           string `json:"name"`
	PhoneNumber    string `json:"phone_number" binding:"required,e164"`
	Password       string `json:"password,omitempty" gorm:"-" binding:"required,min=8,max=15"`
	HashedPassword string `json:"-" gorm:"password"`
	IsEmailActive  bool   `json:"-"`
}
