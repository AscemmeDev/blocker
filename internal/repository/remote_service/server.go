package remoteService

import (
	"blocker/internal/entity"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Server struct {
	client  *http.Client
	BaseURL string
}

func NewServer(baseURL string) *Server {
	return &Server{client: &http.Client{}, BaseURL: baseURL}
}

func (s Server) GetLimits() (n uint64, p time.Duration, err error) {
	var respData GetLimitResponse
	uri := fmt.Sprintf("%s/base/v1/get-limits", s.BaseURL)
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	err = json.Unmarshal(all, &respData)
	if err != nil {
		return 0, 0, err
	}

	return respData.Limit, respData.Duration, nil
}

func (s Server) Process(ctx context.Context, batch entity.Batch) error {
	var requestBatch SendBatch
	var respData SendBatchResponse
	requestBatch.Batch = batch
	marshal, _ := json.Marshal(requestBatch)
	uri := fmt.Sprintf("%s/base/v1/process", s.BaseURL)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewReader(marshal))
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(all, &respData)
	if err != nil {
		return err
	}

	if respData.Code != "200" {
		return errors.New(fmt.Sprintf("status code is:%s, error: %s", respData.Code, respData.Error))
	}

	return nil
}
