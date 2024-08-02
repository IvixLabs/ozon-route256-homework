package gracefulshutdown

import (
	"golang.org/x/net/context"
	"log"
	"route256/common/pkg/env"
)

func GetEndChannel(ctx context.Context) func() {
	endCh := make(chan struct{})

	gracefulShutdownTimeout := env.GetGracefulShutdownTimeout()
	go func() {
		defer close(endCh)

		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
		defer cancel()

		select {
		case <-endCh:
		case <-shutdownCtx.Done():
			log.Fatalln("Shutdown timeout is expired")
		}
	}()

	return func() { endCh <- struct{}{} }
}
