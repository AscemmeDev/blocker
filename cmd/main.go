package main

import (
	"blocker/internal/configs"
	remoteService "blocker/internal/repository/remote_service"
	"blocker/internal/repository/scheduler"
	"blocker/internal/usecases"
	"blocker/internal/worker"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := logrus.New()

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 2)

	newConf := configs.NewConfig()

	newScheduler := scheduler.NewScheduler()
	newServer := remoteService.NewServer(newConf.BaseUrl)

	newRemoteServer := usecases.NewRemoteServer(newServer, logger)
	newBatchSenderUsecase := usecases.NewBatchSenderUsecase(newRemoteServer, ctx, cancel, newScheduler, logger)

	newWorker := worker.NewWorker(newBatchSenderUsecase, logger)

	go func() {
		newWorker.Run(ctx)
		errCh <- errors.New("newWorker is done")

	}()
	go func() {
		newScheduler.Run(ctx)
		errCh <- errors.New("newScheduler is done")
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sig:
		signal.Stop(sig)
		signal.Reset(os.Interrupt)
	case err := <-errCh:
		if err != nil {
			logger.Fatalf("error: %v", err)
		}
	}
}
