package service

import (
	"context"
	"database/sql"

	"github.com/core-go/core/tx"
)

type Repository[T any, K any] interface {
	Load(ctx context.Context, id K) (*T, error)
	Create(ctx context.Context, model *T) (int64, error)
	Update(ctx context.Context, model *T) (int64, error)
	Patch(ctx context.Context, model map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id K) (int64, error)
}

type Service[T any, K any] struct {
	DB         *sql.DB
	Repository Repository[T, K]
	TxKey      string
}

func NewService[T any, K any](db *sql.DB, repository Repository[T, K], opts ...string) *Service[T, K] {
	txKey := "tx"
	if len(opts) > 0 && len(opts[0]) > 0 {
		txKey = opts[0]
	}
	return &Service[T, K]{db, repository, txKey}
}
func (s *Service[T, K]) Load(ctx context.Context, id K) (*T, error) {
	return s.Repository.Load(ctx, id)
}
func (s *Service[T, K]) Create(ctx context.Context, model *T) (int64, error) {
	return tx.ExecuteTx(ctx, s.DB, s.TxKey, func(ctx context.Context) (int64, error) {
		return s.Repository.Create(ctx, model)
	})
}
func (s *Service[T, K]) Update(ctx context.Context, model *T) (int64, error) {
	return tx.ExecuteTx(ctx, s.DB, s.TxKey, func(ctx context.Context) (int64, error) {
		return s.Repository.Update(ctx, model)
	})
}
func (s *Service[T, K]) Patch(ctx context.Context, model map[string]interface{}) (int64, error) {
	return tx.ExecuteTx(ctx, s.DB, s.TxKey, func(ctx context.Context) (int64, error) {
		return s.Repository.Patch(ctx, model)
	})
}
func (s *Service[T, K]) Delete(ctx context.Context, id K) (int64, error) {
	return tx.ExecuteTx(ctx, s.DB, s.TxKey, func(ctx context.Context) (int64, error) {
		return s.Repository.Delete(ctx, id)
	})
}
