package remoteService

import "time"

type GetLimitResponse struct {
	Duration time.Duration `json:"duration"`
	Limit    uint64        `json:"limit"`
}

type SendBatchResponse struct {
	Error string `json:"error,omitempty"`
	Code  string `json:"code,omitempty"`
}
