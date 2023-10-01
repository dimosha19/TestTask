package deamon

import (
	"fmt"
	"net/http"
)

type Option func(s *http.Server)

func WithPort(port string) Option {
	return func(s *http.Server) {
		s.Addr = fmt.Sprintf(":%s", port)
	}
}

func WithHandler(handler http.Handler) Option {
	return func(s *http.Server) {
		s.Handler = handler
	}
}

func NewServer(options ...Option) *http.Server {
	srv := &http.Server{
		Addr: "0.0.0.0:8080",
	}

	for _, opt := range options {
		opt(srv)
	}

	return srv
}
