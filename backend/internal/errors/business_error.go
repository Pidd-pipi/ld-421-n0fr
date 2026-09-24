package errors

import "fmt"

// BusinessError 业务错误，携带统一错误码与 HTTP 状态。
type BusinessError struct {
	Code       int
	Message    string
	HTTPStatus int
	Err        error
}

func (e *BusinessError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *BusinessError) Unwrap() error { return e.Err }

// NewBusinessError 创建业务错误。
func NewBusinessError(code, httpStatus int, message string) *BusinessError {
	return &BusinessError{Code: code, HTTPStatus: httpStatus, Message: message}
}

// WrapBusinessError 包裹底层错误，保留错误链。
func WrapBusinessError(code, httpStatus int, message string, err error) *BusinessError {
	return &BusinessError{Code: code, HTTPStatus: httpStatus, Message: message, Err: err}
}
