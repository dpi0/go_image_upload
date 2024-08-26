package utils

import (
	"github.com/google/uuid"
)

func GenerateShortUUID() string {
	return uuid.New().String()[:5]
}
