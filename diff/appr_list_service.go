package diff

import "context"

type ApprListService interface {
	Approve(ctx context.Context, ids interface{}) (int, error)
	Reject(ctx context.Context, ids interface{}) (int, error)
}
