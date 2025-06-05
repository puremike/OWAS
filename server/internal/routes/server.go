package routes

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func RunServer(mux http.Handler, PORT string, logger *zap.SugaredLogger) error {

	srv := &http.Server{
		Addr:         ":" + PORT,
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

		logger.Infow("Shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
		defer cancel()

		shutdown <- srv.Shutdown(ctx)
	}()

	logger.Infow("Starting server on port:", "port", PORT)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err := <-shutdown; err != nil {
		return err
	}

	return nil
}
