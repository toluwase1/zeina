package models

import "time"

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
