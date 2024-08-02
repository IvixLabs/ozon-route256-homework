package main

import (
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"route256/common/pkg/gracefulshutdown"
	"route256/debugsrv/pkg/debugsrv"
	"route256/metrics/pkg/prometheus"
	"route256/pprof/pkg/pprof"
	"sync"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		debugsrv.Run(ctx,
			prometheus.WithPrometheus(),
			pprof.WithPprof())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runRest(ctx)
	}()

	done := gracefulshutdown.GetEndChannel(ctx)
	wg.Wait()

	done()
}
