package models

import "time"

type Base struct {
	ID *ID `json:"id,omitempty"`
	BaseDate
}

type BaseDate struct {
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	ActiveAt       *time.Time `json:"active_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	VerifiedAt     *time.Time `json:"verified_at,omitempty"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
}
