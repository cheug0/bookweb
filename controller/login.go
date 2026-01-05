package controller

import (
	"bookweb/dao"
	"bookweb/model"
	"bookweb/service"
	"encoding/json"
	"fmt"
	"html/template"
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
				Path:     "/",
			}
			// 发送Cookie
			http.SetCookie(w, &cookie)

			// 返回JSON成功
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": "登录成功",
			})
		} else {
			// 登录失败，用户名或密码错误
			fmt.Println("登录失败：用户名或密码错误")
			// 返回JSON失败
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "用户名或密码错误",
			})
		}
	} else {
		// GET请求，渲染登录页面
		data := GetCommonData(r).Add("CurrentTitle", "用户登录")

		tPath, ok := GetTplPathOrError(w, "login.html")
		if !ok {
			return
		}
		t, err := template.ParseFiles(tPath, TplPath("head.html"), TplPath("foot.html"))
		if err != nil {
			http.Error(w, "解析模板失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, data)
	}
}
