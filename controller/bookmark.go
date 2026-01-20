package controller

import (
	"bookweb/dao"
	"bookweb/service"
	"bookweb/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

// AddBookmark 添加书签
func AddBookmark(w http.ResponseWriter, r *http.Request) {
	isLogin, sess := dao.IsLogin(r)
	w.Header().Set("Content-Type", "application/json")
	if !isLogin {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "未登录"})
		return
	}

	articleIDStr := r.PostFormValue("articleid")
	articleID, _ := strconv.Atoi(articleIDStr)
	// ID 转换
	articleID = utils.DecodeID(articleID)
	chapterIDStr := r.PostFormValue("chapterid")
	chapterID, _ := strconv.Atoi(chapterIDStr)

	if articleID <= 0 || chapterID <= 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "无效的参数"})
		return
	}

	err := service.AddToBookmark(sess.UserID, articleID, chapterID)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "添加成功"})
}

// DeleteBookmark 删除书签
func DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	isLogin, _ := dao.IsLogin(r)
	w.Header().Set("Content-Type", "application/json")
	if !isLogin {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "未登录"})
		return
	}

	bookIDStr := r.PostFormValue("bookid")
	bookID, _ := strconv.Atoi(bookIDStr)

	if bookID <= 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "无效的书签ID"})
		return
	}

	err := service.RemoveFromBookmark(bookID)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "删除成功"})
}
