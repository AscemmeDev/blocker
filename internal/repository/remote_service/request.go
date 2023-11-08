package remoteService

import "blocker/internal/entity"

type SendBatch struct {
	Batch entity.Batch `json:"batch"`
}
