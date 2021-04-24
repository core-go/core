package uuid

import (
	"context"
	"github.com/google/uuid"
	"strings"
)

func RandomId() string {
	id := uuid.New()
	return strings.Replace(id.String(), "-", "", -1)
}

func Generate(ctx context.Context) (string, error) {
	return RandomId(), nil
}
