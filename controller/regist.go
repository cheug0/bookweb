package controller

import (
	"bookweb/service"
	"encoding/json"
	"fmt"
	"html/template"
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

		tPath, ok := GetTplPathOrError(w, "regist.html")
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
