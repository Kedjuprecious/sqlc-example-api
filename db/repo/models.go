// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repo

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Customer struct {
	ID           string           `json:"id"`
	CustomerName string           `json:"customer_name"`
	Phone        string           `json:"phone"`
	Email        string           `json:"email"`
	CreatedAt    pgtype.Timestamp `json:"created_at"`
}

type Order struct {
	ID          string           `json:"id"`
	CustomerID  string           `json:"customer_id"`
	OrderStatus string           `json:"order_status"`
	TotalPrice  string           `json:"total_price"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
}

type Product struct {
	ID           string           `json:"id"`
	ProductName  string           `json:"product_name"`
	Description  string           `json:"description"`
	ProductPrice pgtype.Numeric   `json:"product_price"`
	CreatedAt    pgtype.Timestamp `json:"created_at"`
}
