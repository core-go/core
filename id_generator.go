package service

import "context"

type IdGenerator interface {
	Generate(ctx context.Context, model interface{}) (int, error)
}
