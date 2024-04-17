package utils

import (
	"strings"

	"github.com/google/uuid"
)

func NewUuid32() string {
	uuidObj := uuid.New()
	uuidStr := strings.Replace(uuidObj.String(), "-", "", -1)
	return uuidStr
}
