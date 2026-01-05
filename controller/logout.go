package controller

import (
	"bookweb/dao"
	"net/http"
)

// Logout 处理用户退出请求
func Logout(w http.ResponseWriter, r *http.Request) {
	// 获取Cookie
	cookie, err := r.Cookie("user_session")
	if err != nil {
		// 没有Session，直接重定向到首页
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// 获取SessionID
	sessID := cookie.Value

	// 删除Session
	dao.DeleteSession(sessID)

	// 使Cookie失效
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	// 重定向到首页
	http.Redirect(w, r, "/", http.StatusFound)
}
