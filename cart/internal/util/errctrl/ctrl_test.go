package errctrl

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCtrl_Positive(t *testing.T) {
	t.Parallel()

	errCtrl, errCtx := NewCtrl(context.Background())

	fn1, ok := errCtrl.Wrap(func(ctx context.Context) error {
		return nil
	})
	assert.NotNil(t, fn1)
	assert.True(t, ok)

	fn2, ok := errCtrl.Wrap(func(ctx context.Context) error {
		return nil
	})
	assert.True(t, ok)
	assert.NotNil(t, fn2)

	assert.NotNil(t, errCtrl.ctx)
	assert.NotNil(t, errCtrl.cancel)
	assert.NoError(t, errCtrl.err)

	fn1(errCtx)
	fn2(errCtx)

	err := errCtrl.Wait()
	assert.NoError(t, err)
}

func TestCtrl_Fn_error(t *testing.T) {
	t.Parallel()

	errCtrl, errCtx := NewCtrl(context.Background())

	fn1, ok := errCtrl.Wrap(func(ctx context.Context) error {
		return nil
	})
	assert.NotNil(t, fn1)
	assert.True(t, ok)

	errSome := errors.New("some error")

	fn2, ok := errCtrl.Wrap(func(ctx context.Context) error {
		return errSome
	})
	assert.True(t, ok)
	assert.NotNil(t, fn2)

	assert.NotNil(t, errCtrl.ctx)
	assert.NotNil(t, errCtrl.cancel)
	assert.NoError(t, errCtrl.err)

	fn2(errCtx)
	fn1(errCtx)

	fn3, ok := errCtrl.Wrap(func(ctx context.Context) error {
		return nil
	})
	assert.Nil(t, fn3)
	assert.False(t, ok)

	err := errCtrl.Wait()
	assert.ErrorIs(t, err, errSome)
}
