package service

import "context"

type ApprService interface {
	Approve(ctx context.Context, id interface{}) (Status, error)
	Reject(ctx context.Context, id interface{}) (Status, error)
}
