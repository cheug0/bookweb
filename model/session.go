// session.go
// 会话模型
// 定义用户/管理员会话的数据结构
package model

// Session 用户会话结构体
type Session struct {
	SessionID string
	UserID    int
	Username  string
}
