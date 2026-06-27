package transactions

import "github.com/google/uuid"

// parseUUID is a small wrapper used by DTO conversion. It returns uuid.Nil on
// parse error so the caller can safely pass-through malformed input rather
// than crashing the request. Validation should be done at the binding layer.
func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}