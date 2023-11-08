package worker

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type Worker struct {
	usecases IBatchSender
	logger   *logrus.Logger
}

func NewWorker(usecases IBatchSender, logger *logrus.Logger) *Worker {
	return &Worker{usecases: usecases, logger: logger}
}

func (w *Worker) Run(cxt context.Context) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-cxt.Done():
			logrus.Warn("ctx is done")
			return

		case <-ticker.C:
			data := w.usecases.MockGenerator()
			w.usecases.Process(data)

		}
	}
}
