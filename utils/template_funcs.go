package utils

import (
	"fmt"
	"html/template"
	"time"
)

// GetAdContent 获取广告内容
func GetAdContent(slotID string) template.HTML {
	if GetAdContentFunc != nil {
		return GetAdContentFunc(slotID)
	}
	return ""
}

// BookFuncMap 书籍页面的模版函数（共享给插件使用）
var BookFuncMap = template.FuncMap{
	"formatSize": func(size int) string {
		if size >= 10000 {
			return fmt.Sprintf("%.1f万", float64(size)/10000.0)
		}
		return fmt.Sprintf("%d", size)
	},
	"formatDate": func(t int64) string {
		if t == 0 {
			return "-"
		}
		return time.Unix(t, 0).Format("2006-01-02")
	},
	"safe": func(s string) template.HTML {
		return template.HTML(s)
	},
	"cover": func(id int) string {
		return GetCoverPath(id)
	},
	"bookUrl": func(id int) string {
		return BookUrl(id)
	},
	"bookIndexUrl": func(id int) string {
		return BookIndexUrl(id)
	},
	"bookIndexPageUrl": func(id, page int) string {
		return BookIndexPageUrl(id, page)
	},
	"readUrl": func(aid, cid int) string {
		return ReadUrl(aid, cid)
	},
	"langtailUrl": func(lid int) string {
		return LangtailUrl(lid)
	},
	"ad": func(slotID string) template.HTML {
		return GetAdContent(slotID)
	},
	"transID": func(id int) int {
		return EncodeID(id)
	},
}
