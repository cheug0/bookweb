package controller

import (
	"bookweb/dao"
	"bookweb/model"
	"bookweb/service"
	"encoding/json"
	"fmt"
	"net/http"
)

// UserCenter 用户中心页面
func UserCenter(w http.ResponseWriter, r *http.Request) {
	// 判断是否登录
	isLogin, sess := dao.IsLogin(r)
	if !isLogin {
		// 未登录，重定向到登录页面
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// 获取用户详细信息
	user, err := dao.GetUserByID(sess.UserID)
	if err != nil {
		// Log error?
		fmt.Println("GetUserByID error:", err)
		// Fallback to session username if DB fails?
		// Or show error page?
		// For now, let's just use session username and empty fields if error
		// But better to just handle error gracefully.
	}

	// 使用预编译模板
	t := GetRenderTemplate(w, r, "user_center.html")
	if t == nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	// Prepare data
	data := GetCommonData(r).Add("CurrentTitle", "个人中心")

	if user != nil {
		data.Add("User", user)

		// 获取书架
		bookcases, _, err := service.GetBookcaseList(user.Id, 1, 100)
		if err == nil {
			data.Add("Bookcases", bookcases)
		} else {
			fmt.Println("GetBookcaseList error:", err)
		}

		// 获取书签
		bookmarks, _, err := service.GetBookmarkList(user.Id, 1, 100)
		if err == nil {
			data.Add("Bookmarks", bookmarks)
		} else {
			fmt.Println("GetBookmarkList error:", err)
		}

	} else {
		data.Add("User", &model.User{Username: sess.Username})
	}

	t.Execute(w, data)
}

// UpdateUser 处理更新用户信息请求
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// 判断是否登录
	isLogin, sess := dao.IsLogin(r)
	if !isLogin {
		// 未登录
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method == "POST" {
		password := r.PostFormValue("password")
		email := r.PostFormValue("email")

		// 调用Service层更新信息
		err := service.UpdateUserInfo(sess.UserID, password, email)

		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "修改成功",
		})
	}
}
