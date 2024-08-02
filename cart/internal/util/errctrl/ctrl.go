package errctrl

import (
	"context"
	"sync"
)

type Ctrl struct {
	ctx    context.Context
	err    error
	cancel context.CancelFunc
	wg     sync.WaitGroup
	once   sync.Once
}

func NewCtrl(ctx context.Context) (*Ctrl, context.Context) {
	localCtx, cancel := context.WithCancel(ctx)

	return &Ctrl{
		ctx:    localCtx,
		cancel: cancel,
		wg:     sync.WaitGroup{},
	}, localCtx
}

func (c *Ctrl) Wrap(fn func(ctx context.Context) error) (func(ctx context.Context), bool) {
	ctxErr := c.ctx.Err()
	if ctxErr != nil {
		return nil, false
	}

	return func(_ context.Context) {
			ctxErr := c.ctx.Err()
			if ctxErr != nil {
				return
			}

			c.wg.Add(1)
			defer c.wg.Done()

			err := fn(c.ctx)

			if err != nil {
				c.once.Do(func() {
					c.err = err
					c.cancel()
				})
			}
		},
		true
}

func (c *Ctrl) Wait() error {
	c.wg.Wait()
	err := c.err

	if err != nil {
		return err
	}

	return nil
}
