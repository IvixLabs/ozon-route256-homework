package test

import (
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	commonTransaction "route256/common/pkg/manager/transaction"
	"route256/common/pkg/manager/transaction/mock"
	"route256/loms/internal/manager/transaction"
	"testing"
)

func TestManager_Commit(t *testing.T) {
	t.Parallel()

	mockCtrl := minimock.NewController(t)
	ctrl := mock.NewControllerMock(mockCtrl)

	ctx := context.Background()

	tx := mock.NewTransactionMock(mockCtrl)
	tx.CommitMock.Expect(ctx).Return(nil)

	ctrl.BeginMock.Expect(ctx).Return(tx, ctx, nil)

	manager := transaction.NewManager([]commonTransaction.Controller{ctrl})

	err := manager.Do(ctx, func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
}

func TestManager_Failed_begin(t *testing.T) {
	t.Parallel()

	errBegin := errors.New("begin is failed")

	mockCtrl := minimock.NewController(t)
	ctrl := mock.NewControllerMock(mockCtrl)

	tx := mock.NewTransactionMock(mockCtrl)

	ctx := context.Background()
	ctrl.BeginMock.Expect(ctx).Return(tx, ctx, errBegin)

	manager := transaction.NewManager([]commonTransaction.Controller{ctrl})

	err := manager.Do(ctx, func(ctx context.Context) error {
		return nil
	})

	assert.ErrorIs(t, err, errBegin)
}

func TestManager_Rollback(t *testing.T) {
	t.Parallel()

	errFn := errors.New("fn is failed")

	mockCtrl := minimock.NewController(t)
	ctrl := mock.NewControllerMock(mockCtrl)

	ctx := context.Background()

	tx := mock.NewTransactionMock(mockCtrl)
	tx.RollbackMock.Expect(ctx).Return(nil)

	ctrl.BeginMock.Expect(ctx).Return(tx, ctx, nil)

	manager := transaction.NewManager([]commonTransaction.Controller{ctrl})

	err := manager.Do(ctx, func(ctx context.Context) error {
		return errFn
	})

	assert.ErrorIs(t, err, errFn)
}

func TestManager_Commit_failed(t *testing.T) {
	t.Parallel()

	errCommit := errors.New("commit is failed")

	defer func() {
		if r := recover(); r == nil {
			assert.ErrorIs(t, errCommit, r.(error))
		}
	}()

	mockCtrl := minimock.NewController(t)
	ctrl := mock.NewControllerMock(mockCtrl)

	ctx := context.Background()

	tx := mock.NewTransactionMock(mockCtrl)
	tx.CommitMock.Expect(ctx).Return(errCommit)
	tx.RollbackMock.Expect(ctx).Return(nil)

	ctrl.BeginMock.Expect(ctx).Return(tx, ctx, nil)

	manager := transaction.NewManager([]commonTransaction.Controller{ctrl})

	err := manager.Do(ctx, func(ctx context.Context) error {
		return nil
	})

	assert.ErrorIs(t, err, errCommit)
}

func TestManager_Rollback_failed(t *testing.T) {
	t.Parallel()

	errRollback := errors.New("rollback is failed")

	defer func() {
		if r := recover(); r == nil {
			assert.ErrorIs(t, errRollback, r.(error))
		}
	}()

	mockCtrl := minimock.NewController(t)
	ctrl := mock.NewControllerMock(mockCtrl)

	ctx := context.Background()

	tx := mock.NewTransactionMock(mockCtrl)
	tx.RollbackMock.Expect(ctx).Return(errRollback)

	ctrl.BeginMock.Expect(ctx).Return(tx, ctx, nil)

	manager := transaction.NewManager([]commonTransaction.Controller{ctrl})

	err := manager.Do(ctx, func(ctx context.Context) error {
		return errors.New("some fn error")
	})

	assert.ErrorIs(t, err, errRollback)
}

func TestManager_FailedOneTx(t *testing.T) {
	t.Parallel()

	mockCtrl := minimock.NewController(t)
	ctrl := mock.NewControllerMock(mockCtrl)

	ctx := context.Background()

	tx := mock.NewTransactionMock(mockCtrl)
	tx.RollbackMock.Expect(ctx).Return(nil)

	ctxTx := context.Background()
	ctrl.BeginMock.Expect(ctx).Return(tx, ctxTx, nil)

	ctrl1 := mock.NewControllerMock(mockCtrl)
	tx1 := mock.NewTransactionMock(mockCtrl)
	tx1.RollbackMock.Expect(ctx).Return(nil)

	ctx1 := context.WithValue(ctxTx, "ctxkey", "ctxval")
	ctrl1.BeginMock.Expect(ctx).Return(tx1, ctx1, nil)

	manager := transaction.NewManager([]commonTransaction.Controller{ctrl, ctrl1})

	errSome := errors.New("some error")

	err := manager.Do(ctx, func(ctx context.Context) error {
		if ctx.Value("ctxkey") != nil {
			return errSome
		}
		return nil
	})

	assert.ErrorIs(t, err, errSome)
}
