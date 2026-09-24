package service

// Actor 当前操作者上下文，由鉴权中间件注入。
type Actor struct {
	UserID   uint
	Username string
	Role     string
	IP       string
}
