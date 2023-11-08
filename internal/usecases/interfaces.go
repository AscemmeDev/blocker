package usecases

import (
	"blocker/internal/entity"
	"context"
	"time"
)

type Service interface {
	GetLimits() (n uint64, p time.Duration)
	Process(ctx context.Context, batch entity.Batch) error
}

type IRemoteServer interface {
	GetLimits() (n uint64, p time.Duration, err error)
	Process(ctx context.Context, batch entity.Batch) error
}

type IScheduler interface {
	Run(ctx context.Context)
	Write(key int32, task func(), duration time.Duration)
	Remove(key int32)
	Count() int32
}
