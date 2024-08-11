package model

import "time"

type Product struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Price     float64   `json:"price" db:"price"`
	Stock     int       `json:"stock" db:"stock"`
	Status    int       `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
