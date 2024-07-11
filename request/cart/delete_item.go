package cart

type (
	DeleteItemUri struct {
		CartUuid     string `uri:"cart_uuid"`
		CartItemUuid string `uri:"cart_item_uuid"`
	}
)
