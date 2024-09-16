package http

import (
	"context"
	"ebook/cmd/cronjob/domain"
	"ebook/cmd/cronjob/executor"
	"encoding/json"
	"errors"
	"net/http"
)

var _ executor.Executor = (*Executor)(nil)

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Name() string {
	return "http"
}

func (e *Executor) Exec(ctx context.Context, j domain.CronJob) error {
	type Config struct {
		Endpoint string
		Method   string
	}
	var cfg Config
	err := json.Unmarshal([]byte(j.Cfg), &cfg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(cfg.Method, cfg.Endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		return errors.New("执行失败")
	}
	return nil
}
