package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) server(mux http.Handler) error {

	srv := &http.Server{
		Addr:         ":" + app.config.port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	// implementing graceful shutdown
	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Infow("Shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
		defer cancel()

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("Starting server on port:", "port", app.config.port, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err := <-shutdown; err != nil {
		return err
	}

	return nil
}
