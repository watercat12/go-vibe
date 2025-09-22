package httpserver

import "e-wallet/pkg/config"


type Options func(s *Server) error

func WithConfig(cfg *config.Config) Options {
	return func(s *Server) error {
		s.Config = cfg
		return nil
	}
}
