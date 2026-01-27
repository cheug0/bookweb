// admin.go
// 管理员模型
// 定义后台管理员的数据结构
package model

// Admin 管理员模型
type Admin struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // 存储加密后的密码
}
