package service

import "context"

type ValuesLoader interface {
	Values(ctx context.Context, ids []string) ([]string, error)
}
