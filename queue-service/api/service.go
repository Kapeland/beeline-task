package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"queue-broker/config"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Start(cfg *config.Config) error {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, gCtx := errgroup.WithContext(ctx)

	brokerServer := createBrokerServer(cfg)

	go func() {
		if err := brokerServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			//slog.Error("Failed running gateway server", "description", err.Error())
			cancel()
		}
	}()

	g.Go(func() error {
		<-gCtx.Done()

		if err := brokerServer.Shutdown(gCtx); err != nil {
			//slog.Error("brokerServer.Shutdown", "description", err.Error())
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		//slog.Error("exit reason: " + err.Error())
		return err
	}

	return nil
}
