package model

import "time"

type OrderStatus int

const (
	OrderStatusPending OrderStatus = iota + 1
	OrderStatusFinished
)

type Order struct {
	ID         int         `json:"id" db:"id"`
	Invoice    string      `json:"invoice" db:"invoice"`
	CustomerID int         `json:"customer_id" db:"customer_id"`
	Status     OrderStatus `json:"status" db:"status"`
	CreatedAt  time.Time   `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at,omitempty" db:"updated_at"`

	ProductList []OrderProduct `json:"product_list"`
	TotalAmount float64        `json:"total_amount"`
}

type OrderProduct struct {
	ID         int       `json:"id" db:"id"`
	ProductID  int       `json:"product_id" db:"product_id"`
	Quantity   int       `json:"quantity" db:"quantity"`
	TotalPrice float64   `json:"total_price" db:"total_price"`
	CreatedAt  time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
