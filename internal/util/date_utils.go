package util

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type CustomDate struct {
	time.Time
}

const customLayout = "02-01-2006" // dia-mês-ano

func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		return nil // permite campo vazio
	}

	t, err := time.Parse(customLayout, s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected dd-mm-yyyy: %w", err)
	}
	cd.Time = t
	return nil
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(cd.Time.Format(customLayout))
}
