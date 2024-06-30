package service

import "context"

type Repository[T any, K any] interface {
	Load(ctx context.Context, id K) (*T, error)
	Create(ctx context.Context, model *T) (int64, error)
	Update(ctx context.Context, model *T) (int64, error)
	Patch(ctx context.Context, model map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id K) (int64, error)
}

type Service[T any, K any] struct {
	repository Repository[T, K]
}

func NewService[T any, K any](repository Repository[T, K]) *Service[T, K] {
	return &Service[T, K]{repository}
}
func (s *Service[T, K]) Load(ctx context.Context, id K) (*T, error) {
	return s.repository.Load(ctx, id)
}
func (s *Service[T, K]) Create(ctx context.Context, model *T) (int64, error) {
	return s.repository.Create(ctx, model)
}
func (s *Service[T, K]) Update(ctx context.Context, model *T) (int64, error) {
	return s.repository.Update(ctx, model)
}
func (s *Service[T, K]) Patch(ctx context.Context, model map[string]interface{}) (int64, error) {
	return s.repository.Patch(ctx, model)
}
func (s *Service[T, K]) Delete(ctx context.Context, id K) (int64, error) {
	return s.repository.Delete(ctx, id)
}
