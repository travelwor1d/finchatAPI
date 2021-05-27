package utils

import "github.com/gofrs/uuid"

func NewUUID() string {
	uuidNew, _ := uuid.NewV4()
	return uuidNew.String()
}
