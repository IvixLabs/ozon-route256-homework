package cache

import (
	"errors"
	"golang.org/x/net/context"
	"time"
)

var ErrItemNotFound = errors.New("item is not found")

type Cache[V any] interface {
	Set(context.Context, string, V, time.Duration) error
	Get(context.Context, string) (V, error)
	Delete(context.Context, string) error
}
