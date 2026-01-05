package controller

import (
	"bookweb/config"
	"bookweb/dao"
	"bookweb/model"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"
)

// CommonData 包含全站通用的模版数据
type CommonData map[string]interface{}

// GetCommonData 获取全站通用的模版数据（SEO、用户信息等）
func GetCommonData(r *http.Request) CommonData {
	isLogin, sess := dao.IsLogin(r)
	username := ""
	if isLogin {
		username = sess.Username
	}

	data := CommonData{
		"IsLogin":    isLogin,
		"Username":   username,
		"SiteName":   config.GlobalConfig.Site.SiteName,
		"SiteDomain": config.GlobalConfig.Site.Domain,
		"Analytics":  template.HTML(config.GlobalConfig.Analytics),
	}

	// 计算处理时间
	if startTime, ok := r.Context().Value(model.StartTimeKey).(time.Time); ok {
		// ProcessingComment 改为函数，延迟执行以捕获包含 DB 查询在内的总耗时
		data["ProcessingComment"] = func() template.HTML {
			duration := time.Since(startTime).Seconds()
			return template.HTML(fmt.Sprintf("<!-- Processed in %.6f second(s) -->", duration))
		}
	}

	// 默认应用 index 规则作为基础 SEO
	if rule, ok := config.GlobalConfig.SeoRules["index"]; ok {
		tags := map[string]string{
			"sitename": config.GlobalConfig.Site.SiteName,
			"domain":   config.GlobalConfig.Site.Domain,
		}
		data["CurrentTitle"] = ReplaceSeoTags(rule.Title, tags)
		data["CurrentKeywords"] = ReplaceSeoTags(rule.Keywords, tags)
		data["CurrentDesc"] = ReplaceSeoTags(rule.Description, tags)
	}

	// 动态获取导航栏链接
	sorts, _ := dao.GetAllSorts()
	var sortLinks []map[string]string
	sortRoute := config.GetRouterConfig().GetRoute("sort")
	for _, s := range sorts {
		// 替换参数 :sid 和 :page (默认页码1)
		url := sortRoute
		url = strings.ReplaceAll(url, ":sid", fmt.Sprintf("%d", s.SortID))
		url = strings.ReplaceAll(url, ":page", "1")
		sortLinks = append(sortLinks, map[string]string{
			"Caption": s.Caption,
			"Url":     url,
		})
	}
	data["SortLinks"] = sortLinks

	// 排行榜链接
	topRoute := config.GetRouterConfig().GetRoute("top")
	data["TopUrl"] = topRoute

	return data
}

// Add 方便链式添加数据
func (d CommonData) Add(key string, value interface{}) CommonData {
	d[key] = value
	return d
}

// ApplySeo 根据规则应用 SEO 标签
func (d CommonData) ApplySeo(pageType string, customTags map[string]string) CommonData {
	rule, ok := config.GlobalConfig.SeoRules[pageType]
	if !ok {
		return d
	}

	// 基础标签
	tags := map[string]string{
		"sitename": config.GlobalConfig.Site.SiteName,
		"domain":   config.GlobalConfig.Site.Domain,
	}
	// 合并自定义标签
	for k, v := range customTags {
		tags[k] = v
	}

	d["CurrentTitle"] = ReplaceSeoTags(rule.Title, tags)
	d["CurrentKeywords"] = ReplaceSeoTags(rule.Keywords, tags)
	d["CurrentDesc"] = ReplaceSeoTags(rule.Description, tags)

	return d
}

// ReplaceSeoTags 执行标签替换
func ReplaceSeoTags(tmpl string, tags map[string]string) string {
	res := tmpl
	for tag, val := range tags {
		placeholder := "{" + tag + "}"
		res = strings.ReplaceAll(res, placeholder, val)
	}
	return res
}

// generatePageList 生成简单的页码列表 [1, 2, 3...]
func generatePageList(current, total int) []int {
	var pages []int
	start := current - 5
	if start < 1 {
		start = 1
	}
	end := start + 9
	if end > total {
		end = total
		start = end - 9
		if start < 1 {
			start = 1
		}
	}
	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}
	return pages
}

// TplPath 根据配置加载模版路径，支持 default 回退
func TplPath(name string) string {
	tpl := config.GlobalConfig.Site.Template

	// 1. 尝试从当前主题加载
	path := "template/" + tpl + "/" + name
	if _, err := os.Stat(path); err == nil {
		return path
	}

	// 2. 如果不存在，回退到 default
	path = "template/default/" + name
	if _, err := os.Stat(path); err == nil {
		return path
	}

	// 3. 都不存在
	return ""
}

// GetTplPathOrError 尝试获取模板路径，如果不存在则直接返回 500 错误
func GetTplPathOrError(w http.ResponseWriter, name string) (string, bool) {
	path := TplPath(name)
	if path == "" {
		http.Error(w, "当前模板不存在，请联系管理员修复", http.StatusInternalServerError)
		return "", false
	}
	return path, true
}
