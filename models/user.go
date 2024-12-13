package models

import "time"

type User struct {
	ID           string
	FirstName    string
	MiddleName   string
	LastName     string
	Email        string
	Password     string
	Role         string
	RefreshToken string
	LastLogin    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
