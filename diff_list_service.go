package service

import "context"

type DiffListService interface {
	DiffOfList(ctx context.Context, ids []interface{}) (*[]DiffModel, error)
}
