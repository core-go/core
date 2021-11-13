package v6

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type HealthChecker struct {
	client  *redis.Client
	name    string
	timeout time.Duration
}

func NewRedisHealthChecker(db *redis.Client, name string, timeouts ...time.Duration) *HealthChecker {
	var timeout time.Duration
	if len(timeouts) >= 1 {
		timeout = timeouts[0]
	} else {
		timeout = 4 * time.Second
	}
	return &HealthChecker{client: db, name: name, timeout: timeout}
}

func NewHealthChecker(db *redis.Client, options ...string) *HealthChecker {
	var name string
	if len(options) >= 1 && len(options[0]) > 0 {
		name = options[0]
	} else {
		name = "redis"
	}
	return &HealthChecker{client: db, name: name, timeout: 4 * time.Second}
}

func (s *HealthChecker) Name() string {
	return s.name
}

func (s *HealthChecker) Check(ctx context.Context) (map[string]interface{}, error) {
	cancel := func() {}
	if s.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, s.timeout)
	}
	defer cancel()

	res := make(map[string]interface{})
	checkerChan := make(chan error)
	go func() {
		_, err := s.client.Ping().Result()
		checkerChan <- err
	}()
	select {
	case err := <-checkerChan:
		return res, err
	case <-ctx.Done():
		return res, fmt.Errorf("timeout")
	}
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
