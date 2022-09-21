package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(cfg Config) *zerolog.Logger {

	log := zerolog.New(os.Stdout).Level(zerolog.Level(cfg.Infolevel))

	if cfg.Pretty {
		log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			Level(zerolog.Level(cfg.Infolevel)).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	if cfg.Timestamp {
		log = log.With().Timestamp().Logger()
	}
	if cfg.Caller {
		log = log.With().Caller().Logger()
	}
	return &log
}
