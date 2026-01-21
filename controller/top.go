package controller

import (
	"bookweb/service"
	"log"
	"net/http"
)

// Top 处理排行榜页面请求
func Top(w http.ResponseWriter, r *http.Request) {
	// 获取通用的模版数据并应用 SEO 规则
	data := GetCommonData(r).ApplySeo("top", nil)

	// 获取排行榜数据
	topData, err := service.GetTopData()
	if err != nil {
		http.Error(w, "获取排行榜数据失败", http.StatusInternalServerError)
		return
	}

	data.Add("Top", topData)

	t := GetRenderTemplate(w, r, "top.html")
	if t == nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error (top.html): %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
