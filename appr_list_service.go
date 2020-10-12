package service

import "context"

type ApprListService interface {
	ApproveList(ctx context.Context, ids []interface{}) (Status, error)
	RejectList(ctx context.Context, ids []interface{}) (Status, error)
}
