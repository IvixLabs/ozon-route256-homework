package grpcgw

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"route256/loms/internal/grpcgw/handler"
)

type Server struct {
	config Config
}

func NewServer(config Config) *Server {
	return &Server{config: config}
}

func (s *Server) Run(ctx context.Context) {
	httpConn, err := net.Listen("tcp", s.config.HTTPAddr)
	if err != nil {
		log.Fatalln(err)
	}

	grpcGwHandler, err := handler.NewHandler(s.config.GRPCAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(grpcGwHandler *handler.GrpcGw) {
		closeErr := grpcGwHandler.Close()
		if closeErr != nil {
			log.Panicln(closeErr)
		}
	}(grpcGwHandler)

	httpServer := http.Server{
		Handler: CORSMiddleware(grpcGwHandler),
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
