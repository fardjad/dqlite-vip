package api

import (
	"context"
	"net"
	"net/http"
)

type BackgroundServer interface {
	ListenAndServeInBackground() error
	Shutdown(ctx context.Context) error
}

type BackgroundHttpServer struct {
	*http.Server
}

func (s *BackgroundHttpServer) ListenAndServeInBackground() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	go http.Serve(l, s.Handler)

	return nil
}

type BackgroundServerFactory interface {
	NewServer(addr string, handler http.Handler) BackgroundServer
}

type BackgroundHttpServerFactory struct{}

func (f *BackgroundHttpServerFactory) NewServer(addr string, handler http.Handler) *BackgroundHttpServer {
	return &BackgroundHttpServer{
		&http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}
