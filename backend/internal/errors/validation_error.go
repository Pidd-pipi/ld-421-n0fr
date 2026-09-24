package errors

import "fmt"

// ValidationError 入参校验错误，FieldErrors 记录字段级错误信息。
type ValidationError struct {
	Message     string
	FieldErrors map[string]string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.FieldErrors)
}

// NewValidationError 创建校验错误。
func NewValidationError(message string, fieldErrors map[string]string) *ValidationError {
	return &ValidationError{Message: message, FieldErrors: fieldErrors}
}
