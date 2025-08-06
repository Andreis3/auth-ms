package validator

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	ErrNotBlank  = "this field cannot be blank"
	ErrMaxLength = "cannot be longer than %d characters"
	ErrMinLength = "must be at least %d characters"
)

type Validator struct {
	FieldErrors map[string][]string
}

func New() *Validator {
	return &Validator{
		FieldErrors: make(map[string][]string),
	}
}

func (v *Validator) HasErrors() bool {
	return len(v.FieldErrors) > 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string][]string)
	}

	if slices.Contains(v.FieldErrors[key], message) {
		return
	}

	v.FieldErrors[key] = append(v.FieldErrors[key], message)
}

func (v *Validator) Assert(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func (v *Validator) Merge(other *Validator) {
	if other == nil {
		return
	}

	for key, messages := range other.FieldErrors {
		for _, msg := range messages {
			v.AddFieldError(key, msg)
		}
	}
}

func (v *Validator) Errors() []string {
	errs := make([]string, 0)

	keys := make([]string, 0, len(v.FieldErrors))
	for k := range v.FieldErrors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		for _, msg := range v.FieldErrors[key] {
			errs = append(errs, fmt.Sprintf("%s: %s", key, msg))
		}
	}

	return errs
}

func (v *Validator) FieldErrorsFlat() map[string]any {
	flat := make(map[string]any, len(v.FieldErrors))
	for key, messages := range v.FieldErrors {
		flat[key] = strings.Join(messages, " | ")
	}
	return flat
}

func (v *Validator) FieldErrorsGrouped() map[string]any {
	addressErrors := make(map[int]map[string]any)
	otherErrors := make(map[string]any)

	for field, message := range v.FieldErrorsFlat() {
		if strings.HasPrefix(field, "addresses[") {
			idxStart := strings.Index(field, "[") + 1
			idxEnd := strings.Index(field, "]")
			fieldStart := strings.Index(field, ".") + 1

			indexStr := field[idxStart:idxEnd]
			fieldName := field[fieldStart:]

			index, err := strconv.Atoi(indexStr)
			if err != nil {
				continue
			}

			if addressErrors[index] == nil {
				addressErrors[index] = make(map[string]any)
			}
			addressErrors[index][fieldName] = message
		} else {
			otherErrors[field] = message
		}
	}

	addressList := []map[string]any{}
	for i := range addressErrors {
		if addrErr, ok := addressErrors[i]; ok {
			addressList = append(addressList, addrErr)
		} else {
			addressList = append(addressList, nil)
		}
	}

	result := make(map[string]any)
	maps.Copy(result, otherErrors)

	if len(addressList) > 0 {
		result["addresses"] = addressList
	}

	return result
}

func (v *Validator) Error() string {
	return strings.Join(v.Errors(), "; ")
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}
