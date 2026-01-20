package utils

import (
	"bookweb/config"
	"bookweb/plugin/ads"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	templateCache = make(map[string]*template.Template)
	templateMu    sync.RWMutex
)

// CommonFuncMap 通用模板函数
var CommonFuncMap = template.FuncMap{
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
	"formatDateShort": func(t int64) string {
		if t == 0 {
			return "-"
		}
		return time.Unix(t, 0).Format("01-02")
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
	"readUrl": func(aid, cid int) string {
		return ReadUrl(aid, cid)
	},
	"langtailUrl": func(lid int) string {
		return LangtailUrl(lid)
	},
	"ad": func(slotID string) template.HTML {
		return ads.GetAdContent(slotID)
	},
	"plus":  func(a, b int) int { return a + b },
	"minus": func(a, b int) int { return a - b },
	"transID": func(id int) int {
		return EncodeID(id)
	},
}

// InitTemplates 初始化所有模板（启动时调用）
func InitTemplates() error {
	tpl := config.GlobalConfig.Site.Template
	tplDir := "template/" + tpl

	templates := []struct {
		name  string
		files []string
	}{
		{"index.html", []string{"index.html", "head.html", "foot.html"}},
		{"book_info.html", []string{"book_info.html", "head.html", "foot.html"}},
		{"book_reader.html", []string{"book_reader.html", "head.html", "foot.html"}},
		{"book_list.html", []string{"book_list.html", "head.html", "foot.html"}},
		{"sort.html", []string{"sort.html", "head.html", "foot.html"}},
		{"top.html", []string{"top.html", "head.html", "foot.html"}},
		{"search.html", []string{"search.html", "head.html", "foot.html"}},
		{"user_center.html", []string{"user_center.html", "head.html", "foot.html"}},
		{"login.html", []string{"login.html", "head.html", "foot.html"}},
		{"regist.html", []string{"regist.html", "head.html", "foot.html"}},
		{"error.html", []string{"error.html", "head.html", "foot.html"}},
	}

	templateMu.Lock()
	defer templateMu.Unlock()

	for _, t := range templates {
		var files []string
		for _, f := range t.files {
			path := filepath.Join(tplDir, f)
			if _, err := os.Stat(path); err == nil {
				files = append(files, path)
			} else {
				// 尝试 default 模板
				defaultPath := filepath.Join("template/default", f)
				if _, err := os.Stat(defaultPath); err == nil {
					files = append(files, defaultPath)
				}
			}
		}

		if len(files) > 0 {
			tmpl := template.New(t.name).Funcs(CommonFuncMap)
			tmpl, err := tmpl.ParseFiles(files...)
			if err != nil {
				return fmt.Errorf("error parsing template %s: %v", t.name, err)
			}
			templateCache[t.name] = tmpl
			fmt.Printf("Template cached: %s\n", t.name)
		}
	}

	return nil
}

// GetTemplate 获取预编译的模板
func GetTemplate(name string) *template.Template {
	templateMu.RLock()
	defer templateMu.RUnlock()
	return templateCache[name]
}

// MustGetTemplate 获取模板，如果不存在则 panic
func MustGetTemplate(name string) *template.Template {
	t := GetTemplate(name)
	if t == nil {
		panic("template not found: " + name)
	}
	return t
}
