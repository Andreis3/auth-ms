package util

import "time"

func ToInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func ToInt64Pointer(i int64) *int64 {
	return &i
}

func ToStringPointer(s string) *string {
	return &s
}

func ToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ToTimePointer(t time.Time) *time.Time {
	return &t
}

func ToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
