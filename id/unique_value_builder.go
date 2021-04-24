package id

import "context"

type UniqueValueBuilder interface {
	Build(ctx context.Context, model interface{}, name string) (string, error)
}
