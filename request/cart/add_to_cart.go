package cart

import "github.com/anhhuy1010/customer-cart/models"

type (
	CartRequest struct {
		ProductUuid string `json:"product_uuid" `
		CartUuid    string `json:"cart_uuid"`
	}
	CartResponse struct {
		CartUuid string            `json:"cart_uuid" `
		Items    []models.CartItem `json:"cart_items"`
		Total    int               `json:"total"`
	}
)
