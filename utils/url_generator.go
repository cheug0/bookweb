package utils

import (
	"bookweb/config"
	"strconv"
	"strings"
)

// BookUrl 根据路由配置生成小说信息页 URL
// 从路由 "book" 规则读取模式，替换 :aid 参数
// BookUrl 根据路由配置生成小说信息页 URL
// 从路由 "book" 规则读取模式，替换 :aid 参数
func BookUrl(articleID int) string {
	cfg := config.GetRouterConfig()
	if cfg == nil {
		return "/book_" + strconv.Itoa(EncodeID(articleID)) + ".html"
	}
	pattern := cfg.GetRoute("book")
	if pattern == "" {
		return "/book_" + strconv.Itoa(EncodeID(articleID)) + ".html"
	}
	// 简单的字符串替换比正则快
	return strings.Replace(pattern, ":aid", strconv.Itoa(EncodeID(articleID)), 1)
}

// ReadUrl 根据路由配置生成章节阅读页 URL
// 从路由 "read" 规则读取模式，替换 :aid 和 :cid 参数
func ReadUrl(articleID, chapterID int) string {
	cfg := config.GetRouterConfig()
	if cfg == nil {
		return "/book/" + strconv.Itoa(EncodeID(articleID)) + "/" + strconv.Itoa(chapterID) + "/"
	}
	pattern := cfg.GetRoute("read")
	if pattern == "" {
		// 默认优化路径：直接拼接，避免 Replace 开销
		return "/book/" + strconv.Itoa(EncodeID(articleID)) + "/" + strconv.Itoa(chapterID) + "/"
	}
	// 动态路由替换
	aidStr := strconv.Itoa(EncodeID(articleID))
	cidStr := strconv.Itoa(chapterID)
	// 一次性替换或链式替换
	url := strings.Replace(pattern, ":aid", aidStr, 1)
	return strings.Replace(url, ":cid", cidStr, 1)
}

// GetSiteName 获取网站名称
func GetSiteName() string {
	cfg := config.GetGlobalConfig()
	if cfg != nil && cfg.Site.SiteName != "" {
		return cfg.Site.SiteName
	}
	return "小说网站"
}

// LangtailUrl 根据插件配置生成长尾词页面 URL
func LangtailUrl(langID int) string {
	// 从插件配置获取路由模式
	pluginCfg := config.GetPluginConfig("langtail")
	pattern := "/langtail/:lid" // 默认模式
	if pluginCfg != nil {
		if routePattern, ok := pluginCfg["route_pattern"].(string); ok && routePattern != "" {
			pattern = routePattern
		}
	}
	return strings.Replace(pattern, ":lid", strconv.Itoa(langID), 1)
}

// SortUrl 根据路由配置生成分类列表页 URL
func SortUrl(sortID, page int) string {
	cfg := config.GetRouterConfig()
	if page < 1 {
		page = 1
	}

	defaultUrl := "/sort_" + strconv.Itoa(sortID) + "_" + strconv.Itoa(page) + ".html"

	if cfg == nil {
		return defaultUrl
	}
	pattern := cfg.GetRoute("sort")
	if pattern == "" {
		return defaultUrl
	}

	// 动态路由替换
	sidStr := strconv.Itoa(sortID)
	pageStr := strconv.Itoa(page)

	url := strings.Replace(pattern, ":sid", sidStr, 1)
	return strings.Replace(url, ":page", pageStr, 1)
}
