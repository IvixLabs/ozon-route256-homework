package util

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"log"
	"route256/cart/internal/util/errctrl"
	"route256/cart/internal/util/rppctrl"
	"sync/atomic"
	"testing"
	"time"
)

func TestCtrls_Sequence(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	errCtrl, errCtx := errctrl.NewCtrl(ctx)
	errSome := errors.New("some error")

	var sum atomic.Int32
	gen := func() <-chan func(ctx context.Context) {
		ch := make(chan func(ctx context.Context))

		go func() {
			defer close(ch)

			for i := range 10 {

				wrappedFn, ok := errCtrl.Wrap(
					func(ctx context.Context) error {
						tCtx, cancel := context.WithTimeout(ctx, time.Millisecond)
						defer cancel()

						select {
						case <-tCtx.Done():
							if i == 3 {
								return errSome
							}
							sum.Add(1)
							break
						case <-ctx.Done():
							return ctx.Err()
						}

						return nil
					},
				)

				if !ok {
					break
				}

				ch <- wrappedFn
			}
		}()

		return ch
	}

	rpsCtrl := rppctrl.NewCtrl(errCtx, 2, time.Millisecond*100)

	fns := gen()
	rpsCtrl.GoAll(fns)

	rpsCtrl.Wait()

	err := errCtrl.Wait()

	assert.ErrorIs(t, err, errSome)
	log.Println(sum.Load())
	assert.True(t, sum.Load() < 10)
}
