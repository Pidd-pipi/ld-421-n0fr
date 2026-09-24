package repository

import "errors"

// 仓储层哨兵错误，上层通过 errors.Is 判断。
var (
	ErrNotFound = errors.New("repository: not found")
	ErrConflict = errors.New("repository: conflict")
)
