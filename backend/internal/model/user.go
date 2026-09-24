package model

// Role 系统角色实体。
type Role struct {
	Base
	Code        string `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name        string `gorm:"size:64;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
}

func (Role) TableName() string { return "roles" }

// User 系统用户实体。
type User struct {
	Base
	Username     string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	Name         string `gorm:"size:64;not null" json:"name"`
	Email        string `gorm:"size:128" json:"email"`
	Phone        string `gorm:"size:32" json:"phone"`
	RoleID       uint   `gorm:"index;not null" json:"roleId"`
	Active       bool   `gorm:"default:true" json:"active"`
	Role         *Role  `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (User) TableName() string { return "users" }
