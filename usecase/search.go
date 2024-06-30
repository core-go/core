package usecase

import (
	"context"
	"database/sql"

	"github.com/core-go/core/tx"
)

type SearchRepository[T any, K any, F any] interface {
	Load(ctx context.Context, id K) (*T, error)
	Create(ctx context.Context, model *T) (int64, error)
	Update(ctx context.Context, model *T) (int64, error)
	Patch(ctx context.Context, model map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id K) (int64, error)
	Search(ctx context.Context, filter F, limit int64, offset int64) ([]T, int64, error)
}

type SearchUseCase[T any, K any, F any] struct {
	DB         *sql.DB
	Repository SearchRepository[T, K, F]
	TxKey      string
}

func NewSearchUseCase[T any, K any, F any](db *sql.DB, repository SearchRepository[T, K, F], opts ...string) *SearchUseCase[T, K, F] {
	txKey := "tx"
	if len(opts) > 0 && len(opts[0]) > 0 {
		txKey = opts[0]
	}
	return &SearchUseCase[T, K, F]{db, repository, txKey}
}
func (s *SearchUseCase[T, K, F]) Load(ctx context.Context, id K) (*T, error) {
	return s.Repository.Load(ctx, id)
}
func (s *SearchUseCase[T, K, F]) Create(ctx context.Context, model *T) (int64, error) {
	return tx.ExecuteTx(ctx, s.DB, s.TxKey, func(ctx context.Context) (int64, error) {
		return s.Repository.Create(ctx, model)
	})
}
func (s *SearchUseCase[T, K, F]) Update(ctx context.Context, model *T) (int64, error) {
	return tx.ExecuteTx(ctx, s.DB, s.TxKey, func(ctx context.Context) (int64, error) {
		return s.Repository.Update(ctx, model)
	})
}
func (s *SearchUseCase[T, K, F]) Patch(ctx context.Context, model map[string]interface{}) (int64, error) {
	return tx.ExecuteTx(ctx, s.DB, s.TxKey, func(ctx context.Context) (int64, error) {
		return s.Repository.Patch(ctx, model)
	})
}
func (s *SearchUseCase[T, K, F]) Delete(ctx context.Context, id K) (int64, error) {
	return tx.ExecuteTx(ctx, s.DB, s.TxKey, func(ctx context.Context) (int64, error) {
		return s.Repository.Delete(ctx, id)
	})
}
func (s *SearchUseCase[T, K, F]) Search(ctx context.Context, filter F, limit int64, offset int64) ([]T, int64, error) {
	return s.Repository.Search(ctx, filter, limit, offset)
}
