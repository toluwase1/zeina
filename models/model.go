package models

import "github.com/google/uuid"

type Model struct {
	ID        string `json:"id" gorm:"primaryKey,autoIncrement"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}

func (m *Model) BeforeCreate() error {
	if m.ID == "" {
		m.ID = uuid.NewV4().String()
	}
	return nil
}
