package service

import "context"

type DiffListService interface {
	DiffList(ctx context.Context, ids []interface{}) (*[]DiffModel, error)
}
