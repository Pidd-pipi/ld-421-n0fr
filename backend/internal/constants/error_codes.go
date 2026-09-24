package constants

// 统一响应错误码，0 表示成功，其余为业务/系统错误。
const (
	CodeSuccess      = 0
	CodeBadRequest   = 40000
	CodeUnauthorized = 40100
	CodeForbidden    = 40300
	CodeNotFound     = 40400
	CodeConflict     = 40900
	CodeValidation   = 42200
	CodeInternal     = 50000
)
