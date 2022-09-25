package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Send()
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Err(err).Msg("cant parse config")
	}

	log := NewLogger(cfg)

	ctx, cancel := context.WithCancel(context.Background())

	db := NewDbClient(cfg, log)

	tg := NewTGWeatherAPI(log, cfg, db)

	//Initialize shutdown signal start the long polling the server
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		//shutdown
		db.Disconnect(ctx)
		defer cancel()
		log.Fatal().Msg("user interrupted, exiting")
		os.Exit(1)

	}()

	go tg.PushWeatherUpdates(ctx)

	for {
		for _, upd := range tg.GetUpdates() {
			q, err := tg.ParseMessage(ctx, upd.Message)
			if err != nil {
				tg.log.Error().Err(err).Send()
			}

			_, err = tg.sendMessage(sendMessage, q)
			if err != nil {
				tg.log.Error().Err(err).Send()
			}

			tg.log.Info().Msg(fmt.Sprintf("%+v", upd))
		}
		time.Sleep(time.Second)
	}

}
