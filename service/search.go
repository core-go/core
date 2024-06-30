package service

import "context"

type SearchRepository[T any, K any, F any] interface {
	Load(ctx context.Context, id K) (*T, error)
	Create(ctx context.Context, model *T) (int64, error)
	Update(ctx context.Context, model *T) (int64, error)
	Patch(ctx context.Context, model map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id K) (int64, error)
	Search(ctx context.Context, m F, limit int64, offset int64) ([]T, int64, error)
}

type SearchService[T any, K any, F any] struct {
	repository SearchRepository[T, K, F]
}

func NewSearchService[T any, K any, F any](repository SearchRepository[T, K, F]) *SearchService[T, K, F] {
	return &SearchService[T, K, F]{repository}
}
func (s *SearchService[T, K, F]) Load(ctx context.Context, id K) (*T, error) {
	return s.repository.Load(ctx, id)
}
func (s *SearchService[T, K, F]) Create(ctx context.Context, model *T) (int64, error) {
	return s.repository.Create(ctx, model)
}
func (s *SearchService[T, K, F]) Update(ctx context.Context, model *T) (int64, error) {
	return s.repository.Create(ctx, model)
}
func (s *SearchService[T, K, F]) Patch(ctx context.Context, model map[string]interface{}) (int64, error) {
	return s.repository.Patch(ctx, model)
}
func (s *SearchService[T, K, F]) Delete(ctx context.Context, id K) (int64, error) {
	return s.repository.Delete(ctx, id)
}
func (s *SearchService[T, K, F]) Search(ctx context.Context, filter F, limit int64, offset int64) ([]T, int64, error) {
	return s.repository.Search(ctx, filter, limit, offset)
}
