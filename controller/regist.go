// regist.go
// 注册控制器
// 处理用户注册请求
package controller

import (
	"bookweb/service"
	"encoding/json"
	"fmt"
	"net/http"
)

// Regist 处理用户注册请求
func Regist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// 获取参数
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		email := r.PostFormValue("email")

		// 调用Service层进行注册
		err := service.Register(username, password, email)

		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			fmt.Println("注册失败：", err)
			// 返回JSON错误
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		// 注册成功
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "注册成功，请登录",
		})

	} else {
		// GET请求，渲染注册页面
		data := GetCommonData(r).Add("CurrentTitle", "用户注册")

		// 使用预编译模板（自动根据PC/移动端选择）
		t := GetRenderTemplate(w, r, "regist.html")
		if t == nil {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
		t.Execute(w, data)
	}
}
