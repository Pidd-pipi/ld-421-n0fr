package config

import "fmt"

// DatabaseConfig 数据库连接配置，来源于环境变量。
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// DSN 返回 GORM 使用的 MySQL 连接串。
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}
