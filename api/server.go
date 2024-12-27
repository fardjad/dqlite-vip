package api

import (
	"context"
	"log"
	"net"
	"net/http"
)

type BackgroundServerFactoryFunc func(addr string, handler http.Handler) BackgroundServer

type BackgroundServer interface {
	ListenAndServeInBackground() error
	Shutdown(ctx context.Context) error
}

type BackgroundHttpServer struct {
	*http.Server
}

func NewBackgroundHttpServer(addr string, handler http.Handler) BackgroundServer {
	return &BackgroundHttpServer{
		&http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *BackgroundHttpServer) ListenAndServeInBackground() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	log.Printf("HTTP server is listening on %s", s.Addr)

	go http.Serve(l, s.Handler)

	return nil
}

func (s *BackgroundHttpServer) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
