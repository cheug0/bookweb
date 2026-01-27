// seo_service.go
// SEO 服务
// 处理 SEO 规则的加载和应用
package service

import (
	"bookweb/config"
	"strings"
)

// SeoData 包含 SEO 渲染所需的数据
type SeoData struct {
	Title       string
	Keywords    string
	Description string
}

// GetSeoData 根据页面类型和自定义标签获取 SEO 数据
func GetSeoData(pageType string, customTags map[string]string) SeoData {
	cfg := config.GetGlobalConfig()
	rule, ok := cfg.SeoRules[pageType]
	if !ok {
		// 如果指定的规则不存在，尝试回退到 index 规则
		rule, ok = cfg.SeoRules["index"]
		if !ok {
			return SeoData{}
		}
	}

	// 基础标签
	tags := map[string]string{
		"sitename": cfg.Site.SiteName,
		"domain":   cfg.Site.Domain,
	}
	// 合并自定义标签
	for k, v := range customTags {
		tags[k] = v
	}

	return SeoData{
		Title:       ReplaceSeoTags(rule.Title, tags),
		Keywords:    ReplaceSeoTags(rule.Keywords, tags),
		Description: ReplaceSeoTags(rule.Description, tags),
	}
}

// ReplaceSeoTags 执行 SEO 标签替换逻辑
func ReplaceSeoTags(tmpl string, tags map[string]string) string {
	res := tmpl
	for tag, val := range tags {
		placeholder := "{" + tag + "}"
		res = strings.ReplaceAll(res, placeholder, val)
	}
	return res
}
