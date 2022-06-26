package core

import "context"

type StringService interface {
	Load(ctx context.Context, key string, max int64) ([]string, error)
	Save(ctx context.Context, values []string) (int64, error)
	Delete(ctx context.Context, values []string) (int64, error)
}

type QueryService interface {
	LoadAndDecode(ctx context.Context, id interface{}, result interface{}) (bool, error)
	Search(ctx context.Context, filter interface{}, results interface{}, limit int64, options ...int64) (int64, string, error)
}
type ViewService interface {
	All(ctx context.Context) (interface{}, error)
	Load(ctx context.Context, id interface{}) (interface{}, error)
	LoadAndDecode(ctx context.Context, id interface{}, result interface{}) (bool, error)
	Exist(ctx context.Context, id interface{}) (bool, error)
}
type Service interface {
	ViewService
	Insert(ctx context.Context, model interface{}) (int64, error)
	Update(ctx context.Context, model interface{}) (int64, error)
	Patch(ctx context.Context, model map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id interface{}) (int64, error)
}
type GenericService interface {
	Service
}
