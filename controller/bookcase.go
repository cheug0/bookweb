package controller

import (
	"bookweb/dao"
	"bookweb/service"
	"bookweb/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

// AddBookcase 添加到书架
func AddBookcase(w http.ResponseWriter, r *http.Request) {
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

	if articleID <= 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "无效的文章ID"})
		return
	}

	err := service.AddToBookcase(sess.UserID, articleID)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "添加成功"})
}

// GetBookcaseListAPI 获取书架列表API (JSON)
func GetBookcaseListAPI(w http.ResponseWriter, r *http.Request) {
	isLogin, sess := dao.IsLogin(r)
	w.Header().Set("Content-Type", "application/json")
	if !isLogin {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "未登录"})
		return
	}

	pageStr := r.FormValue("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	list, count, err := service.GetBookcaseList(sess.UserID, page, 20) // Default 20
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    list,
		"count":   count,
	})
}

// DeleteBookcase 删除书架
func DeleteBookcase(w http.ResponseWriter, r *http.Request) {
	isLogin, sess := dao.IsLogin(r)
	w.Header().Set("Content-Type", "application/json")
	if !isLogin {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "未登录"})
		return
	}

	caseIDStr := r.PostFormValue("caseid")
	if caseIDStr != "" {
		caseID, _ := strconv.Atoi(caseIDStr)
		if caseID <= 0 {
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "无效的书架ID"})
			return
		}
		err := service.RemoveBookcaseByID(caseID)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": err.Error()})
			return
		}
	} else {
		articleIDStr := r.PostFormValue("articleid")
		articleID, _ := strconv.Atoi(articleIDStr)
		// ID 转换
		articleID = utils.DecodeID(articleID)
		if articleID <= 0 {
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "无效的文章ID"})
			return
		}
		err := service.RemoveFromBookcase(sess.UserID, articleID)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": err.Error()})
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "删除成功"})
}
