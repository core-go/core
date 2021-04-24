package id

import "context"

type Generator interface {
	Generate(ctx context.Context, name string) (string, error)
	Array(ctx context.Context, name string) ([]string, error)
}
