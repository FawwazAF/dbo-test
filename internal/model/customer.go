package model

import (
	"time"
)

type Customer struct {
	ID             int    `json:"id" db:"id"`
	Username       string `json:"username" db:"username"`
	Password       string `json:"password" db:"-"`
	HashedPassword []byte `db:"password"`

	Name        string    `json:"name" db:"name"`
	Email       string    `json:"email" db:"email"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	DateOfBirth time.Time `json:"date_of_birth" db:"date_of_birth"`
	Address     string    `json:"address" db:"address"`
	Status      int       `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (u *Customer) SetHashedPassword(hashedPassword string) {
	u.HashedPassword = []byte(hashedPassword)
}

func (u *Customer) GetHashedPassword() []byte {
	return u.HashedPassword
}
