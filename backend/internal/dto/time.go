package dto

import (
	"fmt"
	"time"
)

// ParseDate 解析日期字符串，支持 yyyy-MM-dd 与 RFC3339。
func ParseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	layouts := []string{"2006-01-02", time.RFC3339}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date %q", value)
}

// ParseDateTime 解析日期时间字符串，支持 yyyy-MM-dd HH:mm 与 RFC3339。
func ParseDateTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	layouts := []string{"2006-01-02 15:04", "2006-01-02T15:04", "2006-01-02T15:04:05", time.RFC3339}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid datetime %q", value)
}
