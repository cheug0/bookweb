package controller

import (
	"bookweb/dao"
	"bookweb/model"
	"bookweb/service"
	"bookweb/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// Login 处理用户登录请求
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// 获取用户名和密码
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")

		// 调用Service层进行登录验证
		user, err := service.Login(username, password)

		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			fmt.Println("登录失败：", err)
			// 返回JSON错误
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "用户名或密码错误",
			})
			return
		}

		if user != nil {
			// 登录成功
			fmt.Println("登录成功：", user.Username)

			// 生成UUID作为SessionID
			sessID := uuid.New().String()
			// 创建Session对象
			sess := &model.Session{
				SessionID: sessID,
				UserID:    user.Id,
				Username:  user.Username,
			}
			// 保存Session
			dao.AddSession(sess)

			// 创建Cookie
			cookie := http.Cookie{
				Name:     "user_session",
				Value:    sessID,
				HttpOnly: true,
				Secure:   false, // 本地 HTTP 开发环境不能设为 true，否则 Safari 无法保存
				SameSite: http.SameSiteLaxMode,
				Path:     "/",
			}
			// 发送Cookie
			http.SetCookie(w, &cookie)

			// 设置一个非 HttpOnly 的 username cookie 供前端 JS 读取 (用于显示 "欢迎, xxx")
			userCookie := http.Cookie{
				Name:     "username",
				Value:    user.Username,
				HttpOnly: false, // 允许 JS 读取
				Secure:   false,
				Path:     "/",
			}
			http.SetCookie(w, &userCookie)

			// 返回JSON成功
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": "登录成功",
			})
		} else {
			// 登录失败，用户名或密码错误
			utils.LogWarn("Login", "登录失败：用户名或密码错误")
			// 返回JSON失败
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "用户名或密码错误",
			})
		}
	} else {
		// GET请求，渲染登录页面

		// 访问登录页时，为了防止前端缓存的用户名与实际状态不一致，强制清除 username cookie
		userCookie := http.Cookie{
			Name:   "username",
			Value:  "",
			MaxAge: -1,
			Path:   "/",
		}
		http.SetCookie(w, &userCookie)

		data := GetCommonData(r).Add("CurrentTitle", "用户登录")

		// 使用预编译模板（自动根据PC/移动端选择）
		t := GetRenderTemplate(w, r, "login.html")
		if t == nil {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
		t.Execute(w, data)
	}
}
