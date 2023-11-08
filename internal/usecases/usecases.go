package usecases

import (
	"blocker/internal/entity"
	"context"
	"github.com/sirupsen/logrus"
	"math/rand"
)

type BatchSenderUsecase struct {
	service Service
	ctx     context.Context
	cansel  context.CancelFunc

	scheduler IScheduler
	logger    *logrus.Logger
}

func NewBatchSenderUsecase(service Service, ctx context.Context, cansel context.CancelFunc, scheduler IScheduler, logger *logrus.Logger) *BatchSenderUsecase {
	return &BatchSenderUsecase{service: service, ctx: ctx, cansel: cansel, scheduler: scheduler, logger: logger}
}

// Process
// the data will flow asynchronously,
// but the requirements do not say that they should go in turn,
// the main thing is that they get to the maximum
func (u *BatchSenderUsecase) Process(batch []entity.Batch) {
	n, p := u.service.GetLimits()
	if n == 0 && p == 0 {
		return
	}

	if len(batch) < int(n) {
		for i := 0; i < int(n); i++ {
			go func(index int) {
				err := u.service.Process(u.ctx, batch[index]) //I’ll do what’s not critical if one of the batches didn’t sent
				if err != nil {
					u.logger.Error("Process sending batch err:", err)
				}
			}(i)
		}
		return
	}

	batchForDelay := make([]entity.Batch, len(batch)-int(n))
	if n != 0 {
		//I create for convenience and readability
		batchForNow := make([]entity.Batch, n)
		batchForNow = batch[:n]
		batchForDelay = batch[n:]

		for i := 0; i < len(batchForNow); i++ {
			go func(index int) {
				err := u.service.Process(u.ctx, batchForNow[index])
				if err != nil {
					u.logger.Error("Process sending batch err:", err)
				}
			}(i)
		}
	}

	task := func() {
		u.Process(batchForDelay)
	}

	key := u.scheduler.Count() + 1
	u.scheduler.Write(key, task, p)
}

// MockGenerator crates just mock data
func (u *BatchSenderUsecase) MockGenerator() []entity.Batch {
	maxLen := rand.Intn(100)
	batch := make([]entity.Batch, maxLen)
	for i := 0; i < maxLen-1; i++ {
		maxLenItem := rand.Intn(100)
		items := make([]entity.Item, maxLenItem)
		for j := 0; j < maxLenItem-1; j++ {
			items[j] = entity.Item{}
		}
		batch[i] = items
	}
	return batch
}
