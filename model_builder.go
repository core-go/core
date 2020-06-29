package service

import "context"

type ModelBuilder interface {
	BuildToInsert(ctx context.Context, model interface{}) interface{}
	BuildToUpdate(ctx context.Context, model interface{}) interface{}
	BuildToPatch(ctx context.Context, model interface{}) interface{}
	BuildToSave(ctx context.Context, model interface{}) interface{}
}
