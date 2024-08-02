package rppctrl

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestCtrl_Go_one_period(t *testing.T) {
	t.Parallel()

	ctrl := NewCtrl(context.Background(), 10, time.Second)

	beginTime := time.Now()
	var cnt atomic.Int32
	for range 10 {
		ok := ctrl.Go(func(ctx context.Context) {
			cnt.Add(1)
		})
		assert.True(t, ok)
	}
	ctrl.Wait()
	endTime := time.Now()

	diffTime := endTime.Sub(beginTime)

	assert.Equal(t, int32(10), cnt.Load())
	assert.True(t, diffTime < time.Second)
}

func TestCtrl_Go_extended_period(t *testing.T) {
	t.Parallel()

	period := time.Millisecond

	ctrl := NewCtrl(context.Background(), 10, period)

	beginTime := time.Now()
	var cnt atomic.Int32
	for range 15 {
		ok := ctrl.Go(func(ctx context.Context) {
			cnt.Add(1)
		})

		assert.True(t, ok)
	}
	ctrl.Wait()
	endTime := time.Now()

	diffTime := endTime.Sub(beginTime)

	assert.Equal(t, int32(15), cnt.Load())
	assert.True(t, diffTime > period)
}

func TestCtrl_Go_one_period_cancel_ctx(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := NewCtrl(ctx, 10, time.Second)

	beginTime := time.Now()
	var cnt atomic.Int32
	for i := range 10 {
		if i == 5 {
			cancel()
		}

		ok := ctrl.Go(func(ctx context.Context) {
			cnt.Add(1)
		})

		if i >= 5 {
			assert.False(t, ok)
		} else {
			assert.True(t, ok)
		}

	}
	ctrl.Wait()
	endTime := time.Now()

	diffTime := endTime.Sub(beginTime)

	assert.Equal(t, int32(5), cnt.Load())
	assert.True(t, diffTime < time.Second)
}

func TestCtrl_Go_all_one_period(t *testing.T) {
	t.Parallel()

	var cnt atomic.Int32

	producer := func() <-chan func(context.Context) {
		ch := make(chan func(context.Context))

		go func() {
			defer close(ch)

			for range 10 {
				fn := func(ctx context.Context) {
					cnt.Add(1)
				}
				ch <- fn
			}
		}()

		return ch
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := NewCtrl(ctx, 10, time.Millisecond)
	beginTime := time.Now()
	ctrl.GoAll(producer())
	ctrl.Wait()
	endTime := time.Now()

	diffTime := endTime.Sub(beginTime)

	assert.Equal(t, int32(10), cnt.Load())
	assert.True(t, diffTime < time.Millisecond)
}

func TestCtrl_Go_all_extended_period(t *testing.T) {
	t.Parallel()

	var cnt atomic.Int32

	producer := func() <-chan func(context.Context) {
		ch := make(chan func(context.Context))

		go func() {
			defer close(ch)

			for range 15 {
				fn := func(ctx context.Context) {
					cnt.Add(1)
				}
				ch <- fn
			}
		}()

		return ch
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := NewCtrl(ctx, 10, time.Millisecond)
	beginTime := time.Now()
	ctrl.GoAll(producer())
	ctrl.Wait()
	endTime := time.Now()

	diffTime := endTime.Sub(beginTime)

	assert.Equal(t, int32(15), cnt.Load())
	assert.True(t, diffTime > time.Millisecond)
}

func TestCtrl_Go_all_cancel_ctx(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cnt atomic.Int32
	producer := func() <-chan func(context.Context) {
		ch := make(chan func(context.Context))

		go func() {
			defer close(ch)

			for i := range 15 {
				fn := func(ctx context.Context) {
					cnt.Add(1)
				}
				ch <- fn
				if i == 4 {
					cancel()
				}

			}
		}()

		return ch
	}

	ctrl := NewCtrl(ctx, 10, time.Millisecond)
	ctrl.GoAll(producer())
	ctrl.Wait()
	assert.True(t, cnt.Load() < 15)
}

func TestCtrl_Go_all_cancel_ctx_during_fn(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cnt atomic.Int32

	stepCh := make(chan struct{}, 1)
	defer close(stepCh)

	goFnCh := make(chan struct{}, 1)
	defer close(goFnCh)

	producer := func() <-chan func(context.Context) {
		ch := make(chan func(context.Context))

		go func() {
			defer close(ch)

			<-stepCh
			fn := func(ctx context.Context) {
				<-goFnCh
				cnt.Add(1)
			}
			ch <- fn
		}()

		return ch
	}

	ctrl := NewCtrl(ctx, 10, time.Millisecond)
	ctrl.GoAll(producer())
	stepCh <- struct{}{}
	cancel()
	goFnCh <- struct{}{}

	ctrl.Wait()
	assert.Equal(t, int32(0), cnt.Load())
}
