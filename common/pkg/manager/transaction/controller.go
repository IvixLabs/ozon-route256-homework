package transaction

import "context"

type Controller interface {
	Begin(ctx context.Context) (Transaction, context.Context, error)
}

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
