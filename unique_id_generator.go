package service

import "context"

type UniqueIdGenerator interface {
	Generate(ctx context.Context) (string, error)
}
