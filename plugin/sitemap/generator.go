// generator.go (sitemap)
// Sitemap 生成器
// 负责遍历数据库并生成 XML 文件的核心逻辑
package sitemap

import (
	"bookweb/config"
	"bookweb/dao"
	"bookweb/utils"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// URL sitemap URL 结构
type URL struct {
	XMLName    xml.Name `xml:"url"`
	Loc        string   `xml:"loc"`
	LastMod    string   `xml:"lastmod,omitempty"`
	ChangeFreq string   `xml:"changefreq,omitempty"`
	Priority   float64  `xml:"priority,omitempty"`
}

// URLSet sitemap URL 集合
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

// SitemapIndex sitemap 索引结构
type SitemapIndex struct {
	XMLName  xml.Name  `xml:"sitemapindex"`
	XMLNS    string    `xml:"xmlns,attr"`
	Sitemaps []Sitemap `xml:"sitemap"`
}

// Sitemap sitemap 索引中的单个 sitemap
type Sitemap struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

// GenerateSitemap 生成网站地图
func GenerateSitemap() error {
	cfg := GetConfig()
	if cfg == nil || !cfg.Enabled {
		return nil
	}

	// 确保输出目录存在
	outputDir := cfg.OutputPath
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	domain := config.GetGlobalConfig().Site.Domain
	if domain == "" {
		domain = "localhost:8080"
	}
	baseURL := "https://" + domain

	// 获取所有文章
	articles, err := dao.GetAllArticlesForSitemap()
	if err != nil {
		return fmt.Errorf("获取文章列表失败: %v", err)
	}

	utils.LogInfo("Sitemap", "Sitemap: 找到 %d 篇文章", len(articles))

	// 构建 URL 列表
	var urls []URL

	// 添加首页
	urls = append(urls, URL{
		Loc:        baseURL + "/",
		LastMod:    time.Now().Format("2006-01-02"),
		ChangeFreq: "daily",
		Priority:   1.0,
	})

	// 添加文章页面
	if cfg.IncludeBooks {
		for _, article := range articles {
			bookURL := utils.BookUrl(article.ArticleID)
			lastMod := time.Unix(int64(article.LastUpdate), 0).Format("2006-01-02")
			urls = append(urls, URL{
				Loc:        baseURL + bookURL,
				LastMod:    lastMod,
				ChangeFreq: cfg.ChangeFreq,
				Priority:   cfg.Priority,
			})
		}
	}

	// 检查是否需要分割
	if len(urls) > cfg.MaxURLsPerFile {
		return generateSitemapIndex(urls, outputDir, baseURL)
	}

	// 生成单个 sitemap
	return writeSitemapFile(urls, filepath.Join(outputDir, "sitemap.xml"))
}

// generateSitemapIndex 生成 sitemap 索引和分割的 sitemap 文件
func generateSitemapIndex(urls []URL, outputDir, baseURL string) error {
	cfg := GetConfig()
	var sitemaps []Sitemap
	fileCount := (len(urls) + cfg.MaxURLsPerFile - 1) / cfg.MaxURLsPerFile

	for i := 0; i < fileCount; i++ {
		start := i * cfg.MaxURLsPerFile
		end := start + cfg.MaxURLsPerFile
		if end > len(urls) {
			end = len(urls)
		}

		filename := fmt.Sprintf("sitemap-%d.xml", i+1)
		filepath := filepath.Join(outputDir, filename)

		if err := writeSitemapFile(urls[start:end], filepath); err != nil {
			return err
		}

		sitemaps = append(sitemaps, Sitemap{
			Loc:     baseURL + "/sitemap/" + filename,
			LastMod: time.Now().Format("2006-01-02"),
		})
	}

	// 生成索引文件
	index := SitemapIndex{
		XMLNS:    "http://www.sitemaps.org/schemas/sitemap/0.9",
		Sitemaps: sitemaps,
	}

	indexPath := filepath.Join(outputDir, "sitemap.xml")
	return writeXMLFile(index, indexPath)
}

// writeSitemapFile 写入 sitemap 文件
func writeSitemapFile(urls []URL, path string) error {
	urlset := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}
	return writeXMLFile(urlset, path)
}

// writeXMLFile 写入 XML 文件
func writeXMLFile(data interface{}, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 写入 XML 头
	file.WriteString(xml.Header)

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("编码 XML 失败: %v", err)
	}

	utils.LogInfo("Sitemap", "Sitemap: 写入文件 %s", path)
	return nil
}
