package usecases

import (
	"blocker/internal/entity"
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type RemoteServer struct {
	server IRemoteServer
	logger *logrus.Logger
}

func NewRemoteServer(server IRemoteServer, logger *logrus.Logger) *RemoteServer {
	return &RemoteServer{server: server, logger: logger}
}

func (r *RemoteServer) GetLimits() (n uint64, p time.Duration) {
	limits, duration, err := r.server.GetLimits()
	if err != nil {
		r.logger.Error(err)
	}
	return limits, duration
}

func (r *RemoteServer) Process(ctx context.Context, batch entity.Batch) error {
	return r.server.Process(ctx, batch)
}
