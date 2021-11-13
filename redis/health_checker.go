package redis

import (
	"context"
	"github.com/garyburd/redigo/redis"
	"time"
)

type HealthChecker struct {
	Pool    *redis.Pool
	Service string
	Timeout time.Duration
}

func NewRedisHealthChecker(pool *redis.Pool, name string, timeouts ...time.Duration) *HealthChecker {
	var timeout time.Duration
	if len(timeouts) >= 1 {
		timeout = timeouts[0]
	} else {
		timeout = 4 * time.Second
	}
	return &HealthChecker{Pool: pool, Service: name, Timeout: timeout}
}

func NewHealthChecker(pool *redis.Pool, options ...string) *HealthChecker {
	var name string
	if len(options) >= 1 && len(options[0]) > 0 {
		name = options[0]
	} else {
		name = "redis"
	}
	return &HealthChecker{Pool: pool, Service: name, Timeout: 4 * time.Second}
}

func (s *HealthChecker) Name() string {
	return s.Service
}

func (s *HealthChecker) Check(ctx context.Context) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	conn := s.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		return res, err
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
