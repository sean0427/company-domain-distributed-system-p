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

	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`

	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
