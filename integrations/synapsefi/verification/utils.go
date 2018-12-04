package verification

import "github.com/google/uuid"

func generateIdempodencyKey() string {
	return uuid.New().String()
}
