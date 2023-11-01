package domain

import (
	"film-management/config"
)

type OptFunc func(*Opts)

type Opts struct {
	userRepository UserRepository
	cfg            *config.Config
}

func defaultOpts(userRepository UserRepository) Opts {
	return Opts{
		userRepository: userRepository,
	}
}

func WithConfig(cfg *config.Config) OptFunc {
	return func(opts *Opts) {
		opts.cfg = cfg
	}
}
