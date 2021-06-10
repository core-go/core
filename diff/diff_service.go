package diff

import "context"

type DiffService interface {
	Diff(ctx context.Context, id interface{}) (*DiffModel, error)
}
