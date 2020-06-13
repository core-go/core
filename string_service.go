package service

import "context"

type StringService interface {
	Load(ctx context.Context, key string, max int64) ([]string, error)
	Save(ctx context.Context, values []string) (int64, error)
	Delete(ctx context.Context, values []string) (int64, error)
}
