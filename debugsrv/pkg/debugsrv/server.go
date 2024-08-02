package debugsrv

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
)

type Server struct {
	config Config
}

func NewServer(config Config) *Server {
	return &Server{config: config}
}

func (s *Server) Run(ctx context.Context) {
	if len(s.config.options) == 0 {
		return
	}

	httpConn, err := net.Listen("tcp", s.config.HTTPAddr)
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()
	for _, opt := range s.config.options {
		opt(mux)
	}

	httpServer := &http.Server{
		Handler: mux,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		if err = httpServer.Serve(httpConn); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Panicln(err)
			}
		}
	}()

	<-ctx.Done()

	if err := httpServer.Shutdown(context.Background()); err != nil {
		log.Panicln(err)
	}
}
