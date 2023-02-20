package model

import (
	"time"
)

type Company struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Contact string `json:"contact"` // TODO: from user domain

	CreatedBy string    `json:"created_by"`
	Created   time.Time `json:"created"`

	UpdatedBy string    `json:"updated_by"`
	Updated   time.Time `json:"updated"`
}
