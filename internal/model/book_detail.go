package model

import (
	"time"
)

type BookDetail struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	AvailableCopies int        `json:"available_copies"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}
