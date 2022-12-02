package pkg

import "github.com/gofrs/uuid"

// StringToUUID – coverts strings to slice of UUIDs
func StringToUUID(s ...string) ([]uuid.UUID, error) {
	result := make([]uuid.UUID, len(s))
	for i, val := range s {
		uuidValue, err := uuid.FromString(val)
		if err != nil {
			return nil, err
		}
		result[i] = uuidValue
	}
	return result, nil
}
