package service

import (
	"github.com/google/uuid"
	"strings"
)

func RandomId() string {
	id := uuid.New()
	return strings.Replace(id.String(), "-", "", -1)
}
func GenerateId() (string, error) {
	id := RandomId()
	return id, nil
}
