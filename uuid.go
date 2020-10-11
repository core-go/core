package service

import (
	"context"
	"github.com/google/uuid"
	"strings"
)

func RandomId() string {
	id := uuid.New()
	return strings.Replace(id.String(), "-", "", -1)
}
func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}
type UUIDGenerator struct {
}
func (s *UUIDGenerator) Generate(ctx context.Context) (string, error) {
	id := RandomId()
	return id, nil
}
