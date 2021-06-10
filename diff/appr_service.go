package diff

import "context"

type ApprService interface {
	Approve(ctx context.Context, id interface{}) (int, error)
	Reject(ctx context.Context, id interface{}) (int, error)
}
