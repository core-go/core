package builder

import "context"

type UniqueIdGenerator interface {
	Generate(ctx context.Context) (string, error)
}
