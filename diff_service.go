package service

import "context"

type DiffService interface {
	Diff(ctx context.Context, id interface{}) (*DiffModel, error)
}
