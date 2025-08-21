package shared

import "github.com/google/uuid"

type Utils struct{}

func (Utils) UUID() string {

	return uuid.NewString()
}
