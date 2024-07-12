package cart

type (
	CheckoutRequest struct {
		CartUuid     string `json:"cart_uuid"`
		CustomerName string `json:"customer_name"`
		Phone        string `json:"phone"`
		Address      string `json:"address"`
	}
)
