package config

import "time"

// JWTConfig JWT 签发与校验配置。
type JWTConfig struct {
	Secret      string
	ExpireHours int
}

// ExpireDuration 返回 token 有效期。
func (c JWTConfig) ExpireDuration() time.Duration {
	if c.ExpireHours <= 0 {
		return 24 * time.Hour
	}
	return time.Duration(c.ExpireHours) * time.Hour
}
