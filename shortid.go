package service

import (
	"context"
	"github.com/teris-io/shortid"
)

func ShortId() (string, error) {
	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		return "", err
	}
	return sid.Generate()
}

func NewShortIdGenerator() *ShortIdGenerator {
	return &ShortIdGenerator{}
}

type ShortIdGenerator struct {
}
func (s *ShortIdGenerator) Generate(ctx context.Context) (string, error) {
	return ShortId()
}
