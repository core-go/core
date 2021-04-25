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

func Func(auto *bool) func(context.Context) (string, error) {
	if auto != nil && *auto {
		return Generate
	}
	return nil
}

func GetFunc(auto bool) func(context.Context) (string, error) {
	if auto {
		return Generate
	}
	return nil
}
