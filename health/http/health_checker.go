package http

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type HealthChecker struct {
	name    string
	url     string
	timeout time.Duration
}

func NewHealthChecker(name, url string, options ...time.Duration) *HealthChecker {
	if len(options) >= 1 && options[0] > 0 {
		return &HealthChecker{name, url, options[0]}
	} else {
		return &HealthChecker{name, url, 4 * time.Second}
	}
}

func (s *HealthChecker) Name() string {
	return s.name
}

func (s *HealthChecker) Check(ctx context.Context) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	client := http.Client{
		Timeout: s.timeout,
		// never follow redirects
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(s.url)
	if e, ok := err.(net.Error); ok && e.Timeout() {
		return res, fmt.Errorf("time out: %w", e)
	} else if err != nil {
		return res, err
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	if resp.StatusCode >= 500 {
		return res, fmt.Errorf("status code is: %d", resp.StatusCode)
	}
	return res, nil
}

func (s *HealthChecker) Build(ctx context.Context, data map[string]interface{}, err error) map[string]interface{} {
	if err == nil {
		return data
	}
	if data == nil {
		data = make(map[string]interface{}, 0)
	}
	data["error"] = err.Error()
	return data
}
