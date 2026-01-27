// auth.go
// 后台认证控制器
// 处理管理员登录验证及 Session 管理
package admin

import (
	"bookweb/config"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

// AdminSession 管理员会话
type AdminSession struct {
	ID        string
	AdminID   int
	Username  string
	CreatedAt time.Time
}

var (
	adminSessions = make(map[string]*AdminSession)
	sessionLock   sync.RWMutex
)

const (
	AdminSessionCookieName = "admin_session"
	AdminSessionExpiry     = 24 * time.Hour
)

// CreateAdminSession 创建管理员会话
func CreateAdminSession(adminID int, username string) *AdminSession {
	sessionLock.Lock()
	defer sessionLock.Unlock()

	session := &AdminSession{
		ID:        uuid.New().String(),
		AdminID:   adminID,
		Username:  username,
		CreatedAt: time.Now(),
	}
	adminSessions[session.ID] = session
	return session
}

// GetAdminSession 获取会话
func GetAdminSession(sessionID string) *AdminSession {
	sessionLock.RLock()
	defer sessionLock.RUnlock()

	session, ok := adminSessions[sessionID]
	if !ok {
		return nil
	}
	// 检查是否过期
	if time.Since(session.CreatedAt) > AdminSessionExpiry {
		delete(adminSessions, sessionID)
		return nil
	}
	return session
}

// DeleteAdminSession 删除会话
func DeleteAdminSession(sessionID string) {
	sessionLock.Lock()
	defer sessionLock.Unlock()
	delete(adminSessions, sessionID)
}

// IsAdminLoggedIn 检查管理员是否已登录
func IsAdminLoggedIn(r *http.Request) (*AdminSession, bool) {
	cookie, err := r.Cookie(AdminSessionCookieName)
	if err != nil {
		return nil, false
	}
	session := GetAdminSession(cookie.Value)
	if session == nil {
		return nil, false
	}
	return session, true
}

// SetAdminSessionCookie 设置会话 Cookie
func SetAdminSessionCookie(w http.ResponseWriter, session *AdminSession) {
	adminPath := "/admin"
	if config.GlobalConfig != nil && config.GlobalConfig.Site.AdminPath != "" {
		adminPath = config.GlobalConfig.Site.AdminPath
	}

	http.SetCookie(w, &http.Cookie{
		Name:     AdminSessionCookieName,
		Value:    session.ID,
		Path:     adminPath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(AdminSessionExpiry.Seconds()),
	})
}

// ClearAdminSessionCookie 清除会话 Cookie
func ClearAdminSessionCookie(w http.ResponseWriter) {
	adminPath := "/admin"
	if config.GlobalConfig != nil && config.GlobalConfig.Site.AdminPath != "" {
		adminPath = config.GlobalConfig.Site.AdminPath
	}

	http.SetCookie(w, &http.Cookie{
		Name:   AdminSessionCookieName,
		Value:  "",
		Path:   adminPath,
		MaxAge: -1,
	})
}

// AuthMiddleware 认证中间件
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, loggedIn := IsAdminLoggedIn(r)
		if !loggedIn {
			adminPath := "/admin"
			if config.GlobalConfig != nil && config.GlobalConfig.Site.AdminPath != "" {
				adminPath = config.GlobalConfig.Site.AdminPath
			}
			http.Redirect(w, r, adminPath+"/login", http.StatusFound)
			return
		}
		next(w, r)
	}
}
