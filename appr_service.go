package service

import "context"

type ApprService interface {
	Approve(ctx context.Context, id interface{}) (StatusCode, error)
	Reject(ctx context.Context, id interface{}) (StatusCode, error)
}
