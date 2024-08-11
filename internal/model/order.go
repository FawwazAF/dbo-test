package model

import "time"

type Order struct {
	ID          int            `json:"id" db:"id"`
	CustomerID  int            `json:"customer_id" db:"customer_id"`
	TotalAmount float64        `json:"total_amount" db:"total_amount"`
	ProductList []ProductOrder `json:"product_list" db:"product_list"`
	Status      int            `json:"status" db:"status"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
}

type ProductOrder struct {
	ProductID   int     `json:"product_id" db:"id"`
	ProductName string  `json:"product_name" db:"name"`
	Quantity    int     `json:"quantity" db:"quantity"`
	Price       float64 `json:"price" db:"price"`
	Status      int     `json:"status" db:"status"`
}
