package models

import "github.com/google/uuid"

type BaseModel struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}

func (m *BaseModel) BeforeCreate() error {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return nil
}
