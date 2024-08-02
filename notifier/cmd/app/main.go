package main

import (
	"context"
	"os"
	"os/signal"
	"route256/common/pkg/gracefulshutdown"
	"route256/debugsrv/pkg/debugsrv"
	"route256/metrics/pkg/prometheus"
	"route256/pprof/pkg/pprof"
	"sync"
)

func main() {

	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		runConsumer(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		debugsrv.Run(ctx,
			prometheus.WithPrometheus(),
			pprof.WithPprof())
	}()

	done := gracefulshutdown.GetEndChannel(ctx)

	wg.Wait()
	done()
}
