// session_dao.go
// 会话 DAO
// 简单的内存会话管理实现
package dao

import (
	"bookweb/model"
	"errors"
	"net/http"
)

// SessionMap 全局会话存储（内存Map）
var SessionMap = make(map[string]*model.Session)

// AddSession 添加会话
func AddSession(sess *model.Session) {
	SessionMap[sess.SessionID] = sess
}

// DeleteSession 删除会话
func DeleteSession(sessID string) {
	delete(SessionMap, sessID)
}

// GetSession 获取会话
func GetSession(sessID string) (*model.Session, error) {
	if sess, ok := SessionMap[sessID]; ok {
		return sess, nil
	}
	return nil, errors.New("会话不存在")
}

// IsLogin 判断用户是否登录
// 返回：(是否登录, Session对象)
func IsLogin(r *http.Request) (bool, *model.Session) {
	// 获取Cookie
	cookie, err := r.Cookie("user_session")
	if err != nil {
		return false, nil
	}

	// 获取SessionID
	sessionID := cookie.Value

	// 查询Session
	sess, err := GetSession(sessionID)
	if err != nil {
		return false, nil
	}

	return true, sess
}
