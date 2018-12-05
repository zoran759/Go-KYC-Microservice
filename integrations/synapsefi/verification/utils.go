package verification

import "github.com/google/uuid"

// newIdempotencyKey returns new idempotency key.
func newIdempotencyKey() string {
	return uuid.New().String()
}
