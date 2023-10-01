package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"task/internal/deamon"

	"github.com/rs/zerolog/log"
)

func main() {
	srv := deamon.NewServer(
		deamon.WithHandler(deamon.NewEngine()),
	)

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Error().
				Err(err).
				Msg("Error occurred during HTTP server Shutdown")
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Error().
			Err(err).
			Msg("Error occurred during ListenAndServe")
	}
	<-idleConnsClosed
}
