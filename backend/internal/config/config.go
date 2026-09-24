package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config 集中管理所有环境变量配置。
type Config struct {
	AppEnv         string   `env:"APP_ENV" envDefault:"development"`
	ServerPort     string   `env:"SERVER_PORT" envDefault:"8080"`
	DBHost         string   `env:"DB_HOST" envDefault:"127.0.0.1"`
	DBPort         string   `env:"DB_PORT" envDefault:"3306"`
	DBName         string   `env:"DB_NAME" envDefault:"lab_equipment"`
	DBUser         string   `env:"DB_USER" envDefault:"lab_user"`
	DBPassword     string   `env:"DB_PASSWORD" envDefault:"lab_password"`
	JWTSecret      string   `env:"JWT_SECRET" envDefault:"lab-equipment-dev-secret"`
	JWTExpireHours int      `env:"JWT_EXPIRE_HOURS" envDefault:"24"`
	CORSOrigins    []string `env:"CORS_ORIGINS" envSeparator:"," envDefault:"http://localhost:18801"`
}

// Load 从环境变量加载配置。
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

// Database 返回数据库配置视图。
func (c *Config) Database() DatabaseConfig {
	return DatabaseConfig{
		Host:     c.DBHost,
		Port:     c.DBPort,
		Name:     c.DBName,
		User:     c.DBUser,
		Password: c.DBPassword,
	}
}

// JWT 返回 JWT 配置视图。
func (c *Config) JWT() JWTConfig {
	return JWTConfig{Secret: c.JWTSecret, ExpireHours: c.JWTExpireHours}
}
