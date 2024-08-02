package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCart(t *testing.T) {
	t.Parallel()
	cart := NewCart(111)
	assert.Equal(t, UserID(111), cart.UserId)
	assert.Equal(t, map[Sku]CartItem{}, cart.Items)
}

func TestCart_AddItem(t *testing.T) {
	t.Parallel()
	cart := NewCart(111)
	cart.AddItem(CartItem{Sku: 222, Count: 2})
	assert.Equal(t, map[Sku]CartItem{222: {Sku: 222, Count: 2}}, cart.Items)
}
