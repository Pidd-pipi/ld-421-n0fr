package dto

import "time"

// LoginRequest 登录请求。
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 用户信息返回。
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	RoleID   uint   `json:"roleId"`
	RoleCode string `json:"roleCode"`
	RoleName string `json:"roleName"`
}

// LoginResponse 登录返回。
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// DateOnly 用于统一输出日期格式的辅助类型。
type DateOnly time.Time

// MarshalJSON 输出 yyyy-MM-dd。
func (d DateOnly) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	if t.IsZero() {
		return []byte(`null`), nil
	}
	return []byte(`"` + t.Format("2006-01-02") + `"`), nil
}
