package transaction

import (
	"context"
	"route256/common/pkg/manager/transaction"
)

type NoopTransaction struct {
}

func NewNoopTransaction() *NoopTransaction {
	return &NoopTransaction{}
}

func (m *NoopTransaction) Commit(_ context.Context) error {

	return nil
}

func (m *NoopTransaction) Rollback(_ context.Context) error {
	return nil
}

type NoopController struct {
}

func NewNoopController() *NoopController {
	return &NoopController{}
}

func (c *NoopController) Begin(ctx context.Context) (transaction.Transaction, context.Context, error) {

	return NewNoopTransaction(), ctx, nil
}
