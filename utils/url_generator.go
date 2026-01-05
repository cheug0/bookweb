package utils

import (
	"bookweb/config"
	"strconv"
	"strings"
)

// BookUrl 根据路由配置生成小说信息页 URL
// 从路由 "book" 规则读取模式，替换 :aid 参数
func BookUrl(articleID int) string {
	cfg := config.GetRouterConfig()
	if cfg == nil {
		return "/book_" + strconv.Itoa(articleID) + ".html"
	}
	pattern := cfg.GetRoute("book")
	if pattern == "" {
		return "/book_" + strconv.Itoa(articleID) + ".html"
	}
	return strings.Replace(pattern, ":aid", strconv.Itoa(articleID), 1)
}

// ReadUrl 根据路由配置生成章节阅读页 URL
// 从路由 "read" 规则读取模式，替换 :aid 和 :cid 参数
func ReadUrl(articleID, chapterID int) string {
	cfg := config.GetRouterConfig()
	if cfg == nil {
		return "/book/" + strconv.Itoa(articleID) + "/" + strconv.Itoa(chapterID) + "/"
	}
	pattern := cfg.GetRoute("read")
	if pattern == "" {
		return "/book/" + strconv.Itoa(articleID) + "/" + strconv.Itoa(chapterID) + "/"
	}
	url := strings.Replace(pattern, ":aid", strconv.Itoa(articleID), 1)
	url = strings.Replace(url, ":cid", strconv.Itoa(chapterID), 1)
	return url
}
