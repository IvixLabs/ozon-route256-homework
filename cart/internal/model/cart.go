package model

type Cart struct {
	UserId UserID
	Items  map[Sku]CartItem
}

func NewCart(userId UserID) *Cart {
	return &Cart{
		UserId: userId,
		Items:  make(map[Sku]CartItem),
	}
}

func (c *Cart) AddItem(item CartItem) {
	c.Items[item.Sku] = item
}
