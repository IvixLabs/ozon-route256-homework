package rppctrl

import (
	"context"
	"sync"
	"time"
)

type Ctrl struct {
	periodCnt  int
	period     time.Duration
	beginTime  time.Time
	currentCnt int
	goWg       sync.WaitGroup
	goAllWg    sync.WaitGroup
	ctx        context.Context
}

func NewCtrl(ctx context.Context, cnt int, period time.Duration) *Ctrl {
	return &Ctrl{
		ctx:       ctx,
		periodCnt: cnt,
		period:    period,
		beginTime: time.Now(),
	}
}

func (c *Ctrl) Go(fn func(ctx context.Context)) bool {
	if c.ctx.Err() != nil {
		return false
	}

	c.sleep()

	c.goWg.Add(1)

	go func() {
		defer c.goWg.Done()
		fn(c.ctx)
	}()

	return true
}

func (c *Ctrl) sleep() {
	if c.currentCnt == c.periodCnt {
		c.currentCnt = 0
		now := time.Now()
		diff := now.Sub(c.beginTime)

		if diff < c.period {
			remainsTime := c.period - diff
			time.Sleep(remainsTime)
		}

		c.beginTime = time.Now()
	}

	c.currentCnt++
}

func (c *Ctrl) GoAll(fns <-chan func(ctx context.Context)) {
	c.goAllWg.Add(1)
	go func() {
		defer c.goAllWg.Done()

	ForLoop:
		for {
			select {
			case <-c.ctx.Done():
				break ForLoop
			case fn, ok := <-fns:
				if !ok {
					break ForLoop
				}

				okGo := c.Go(fn)
				if !okGo {
					break ForLoop
				}
			}
		}

	}()
}

func (c *Ctrl) Wait() {
	c.goAllWg.Wait()
	c.goWg.Wait()
}
