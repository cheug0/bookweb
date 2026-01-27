// user.go
// 用户模型
// 定义前台用户的数据结构
package model

// User 用户
type User struct {
	Id               int
	Username         string
	Password         string
	Email            string
	LastLoginTime    string
	CurrentLoginTime string
}
