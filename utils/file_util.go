// file_util.go
// 文件工具
// 提供文件读写、路径处理等辅助功能
package utils

import (
	"bookweb/config"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GetChapterFileContent 读取章节内容
// path format: /files/article/txt/{articleID/1000}/{articleID}/{chapterOrder_or_chapterID}.txt
func GetChapterFileContent(articleID, chapterID, chapterOrder int) (string, error) {

	fileName := chapterID
	if chapterOrder > 0 {
		fileName = chapterOrder
	}

	// Construct the relative path (rules remain unchanged)
	subDir1 := articleID / 1000
	relPath := fmt.Sprintf("article/txt/%d/%d/%d.txt", subDir1, articleID, fileName)

	// Read content from configured storage
	data, err := GetFileContent(relPath)
	if err != nil {
		return "章节内容不存在", err
	}

	// 1. GBK to UTF-8 Conversion
	utf8Data, err := GbkToUtf8(data)
	if err != nil {
		// If conversion fails, fallback to original data string
		return string(data), nil
	}

	// 2. Format plain text to HTML
	content := string(utf8Data)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	paras := strings.Split(content, "\n")
	var htmlContent strings.Builder
	for _, p := range paras {
		p = strings.TrimSpace(p)
		if p != "" {
			htmlContent.WriteString("<p>")
			htmlContent.WriteString(p)
			htmlContent.WriteString("</p>")
		}
	}

	return htmlContent.String(), nil
}

// GetFileContent 根据配置从不同存储后端读取文件内容
func GetFileContent(relPath string) ([]byte, error) {
	cfg := config.GetGlobalConfig()
	// 默认使用本地存储 (NFS 模式通常也映射为本地路径)

	var storageType string
	if cfg == nil {
		storageType = "local"
	} else {
		storageType = cfg.Storage.Type
	}

	switch storageType {
	case "oss":
		return readOssData(cfg, relPath)
	case "local":
		fullPath := filepath.Join(cfg.Storage.Local.Path, relPath)
		return os.ReadFile(fullPath)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

// readOssData 从对象存储读取数据
func readOssData(cfg *config.AppConfig, relPath string) ([]byte, error) {
	// 如果设置了自定义域名，优先使用域名拼接
	url := ""
	if cfg.Storage.Oss.Domain != "" {
		url = fmt.Sprintf("%s/%s", strings.TrimRight(cfg.Storage.Oss.Domain, "/"), relPath)
	} else {
		// 默认模式: https://{bucket}.{endpoint}/{path}
		endpoint := strings.TrimLeft(cfg.Storage.Oss.Endpoint, "https://")
		endpoint = strings.TrimLeft(endpoint, "http://")
		url = fmt.Sprintf("https://%s.%s/%s", cfg.Storage.Oss.Bucket, endpoint, relPath)
	}

	// 简单实现：通过 HTTP GET 获取 (适用于公共读或已授权路径)
	// 对于私有存储，后续可在此处引入各个厂商的 SDK (如阿里云 OSS, 腾讯云 COS, AWS S3 等)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("oss storage read failed, status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// GbkToUtf8 转换 GBK 字节流为 UTF-8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// GetCoverPath 获取小说封面图片解析路径 (前端调用)
func GetCoverPath(articleID int) string {
	cfg := config.GetGlobalConfig()
	// 如果使用 OSS 且配置了域名，直接返回 OSS 上的 URL
	if cfg != nil && cfg.Storage.Type == "oss" && cfg.Storage.Oss.Domain != "" {
		subDir := articleID / 1000
		return fmt.Sprintf("%s/article/image/%d/%d/%ds.jpg",
			strings.TrimRight(cfg.Storage.Oss.Domain, "/"), subDir, articleID, articleID)
	}
	// 默认本地模式，返回内部路由路径
	return fmt.Sprintf("/img/%d.jpg", EncodeID(articleID))
}

// GetPhysicalCoverPath 获取小说封面图片在存储中的相对物理路径
// 规则: article/image/{articleID/1000}/{articleID}/{articleID}s.jpg
func GetPhysicalCoverPath(articleID int) string {
	subDir := articleID / 1000
	return fmt.Sprintf("article/image/%d/%d/%ds.jpg", subDir, articleID, articleID)
}

// FormatDate 格式化时间戳
func FormatDate(t int64, layout string) string {
	if t == 0 {
		return "-"
	}
	return time.Unix(t, 0).Format(layout)
}
