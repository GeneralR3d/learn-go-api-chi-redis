package model

import (
	"github.com/google/uuid"	//	provides an uuid types
)

type Order struct{
	//	type UUIDv4 is better, longer, more random, less collisions but not very user friendly if we want to display this ID to the user
	OrderID uint64	`json: "order_id"`
	CustomerID uuid.UUID	`json: "customer_id"`
	LineItems []LineItem	`json: "line_items"`
	//OrderStatus can be better represented as timestamps
	CreatedAt *time.Time	`json: "created_at"`
	ShippedAt *time.Time	`json: "shipped_at"`
	CompletedAt *time.Time	`json: "completed_at"`

}

type LineItem struct {
	ItemID uuid.UUID	`json: "item_id"`
	Quantity uint	`json: "quantity"`
	Price uint	`json: "price"`	//	price they paid at the time or purchase
}