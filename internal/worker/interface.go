package worker

import "blocker/internal/entity"

type IBatchSender interface {
	Process(batch []entity.Batch)
	MockGenerator() []entity.Batch
}
