package cart

import "time"

type (
	GetCartRequestUri struct {
		CartUuid string `uri:"cart_uuid"`
	}
	GetCartResponse struct {
		CartUuid string                `json:"cart_uuid"`
		Total    float64               `json:"total"`
		Items    []GetCartItemResponse `json:"items"`
	}
	GetCartItemResponse struct {
		ProductUuid  string    `json:"product_uuid"`
		CartUuid     string    `json:"cart_uuid"`
		CartItemUuid string    `json:"cart_item_uuid"`
		ProductName  string    `json:"product_name"`
		ProductPrice float64   `json:"product_price"`
		Quantity     int       `json:"quantity"`
		ProductTotal float64   `json:"product_total"`
		IsActive     int       `json:"is_active"`
		IsDelete     int       `json:"is_delete"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}
)
