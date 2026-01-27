// template.go
// 模版工具
// HTML 模版的加载、解析与渲染管理
package utils

import (
	"bookweb/config"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sync"
)

var (
	templateCache    = make(map[string]*template.Template)
	templateMu       sync.RWMutex
	GetAdContentFunc func(slotID string) template.HTML
)

// InitTemplates 初始化所有模板（启动时调用）
func InitTemplates() error {
	tpl := config.GlobalConfig.Site.Template
	mobileTpl := config.GlobalConfig.Site.MobileTemplate
	tplDir := "template/" + tpl
	mobileTplDir := "template/" + mobileTpl

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

	// 清空旧缓存（避免 Reload 时残留）
	templateCache = make(map[string]*template.Template)

	for _, t := range templates {
		// 1. 加载 PC 模板
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
				return fmt.Errorf("error parsing PC template %s: %v", t.name, err)
			}
			templateCache[t.name] = tmpl
			LogDebug("Template", "PC Template cached: %s", t.name)
		}

		// 2. 加载移动端模板 (如果配置了)
		if mobileTpl != "" {
			var mFiles []string
			for _, f := range t.files {
				// 优先从移动端模板目录加载
				path := filepath.Join(mobileTplDir, f)
				if _, err := os.Stat(path); err == nil {
					mFiles = append(mFiles, path)
				} else {
					// 直接尝试 default 目录 (不再回退到 PC 模板)
					defaultPath := filepath.Join("template/default", f)
					if _, err := os.Stat(defaultPath); err == nil {
						mFiles = append(mFiles, defaultPath)
					}
				}
			}

			if len(mFiles) > 0 {
				mTmpl := template.New(t.name).Funcs(CommonFuncMap)
				mTmpl, err := mTmpl.ParseFiles(mFiles...)
				if err != nil {
					// 移动端模板加载失败不应该阻断启动，打日志即可
					LogWarn("Template", "Error parsing Mobile template %s: %v", t.name, err)
				} else {
					templateCache["mobile/"+t.name] = mTmpl
					LogDebug("Template", "Mobile Template cached: mobile/%s (files: %v)", t.name, mFiles)
				}
			}
		}
	}

	return nil
}

// GetTemplate 获取 PC 预编译模板
func GetTemplate(name string) *template.Template {
	templateMu.RLock()
	defer templateMu.RUnlock()
	return templateCache[name]
}

// GetMobileTemplate 获取移动端预编译模板
func GetMobileTemplate(name string) *template.Template {
	templateMu.RLock()
	defer templateMu.RUnlock()
	if t, ok := templateCache["mobile/"+name]; ok {
		return t
	}
	// 如果没有移动端专属缓存，回退到 PC 模板
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
