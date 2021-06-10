package diff

import "context"

type DiffListService interface {
	Diff(ctx context.Context, ids interface{}) (*[]DiffModel, error)
}
