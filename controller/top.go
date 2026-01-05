package controller

import (
	"bookweb/service"
	"bookweb/utils"
	"html/template"
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

	// 定义模版函数
	funcMap := template.FuncMap{
		"plus": func(a, b int) int { return a + b },
		"cover": func(id int) string {
			return utils.GetCoverPath(id)
		},
		"bookUrl": func(id int) string {
			return utils.BookUrl(id)
		},
		"readUrl": func(aid, cid int) string {
			return utils.ReadUrl(aid, cid)
		},
	}

	tPath, ok := GetTplPathOrError(w, "top.html")
	if !ok {
		return
	}
	t := template.New("top.html").Funcs(funcMap)
	t, err = t.ParseFiles(tPath, TplPath("head.html"), TplPath("foot.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error (top.html): %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
