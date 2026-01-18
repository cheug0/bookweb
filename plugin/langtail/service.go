package langtail

import (
	"bookweb/dao"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// FetchBaiduSuggestions 从百度搜索建议接口获取长尾词
func FetchBaiduSuggestions(keyword string) ([]string, error) {
	// 构建请求URL
	apiURL := fmt.Sprintf("http://suggestion.baidu.com/su?wd=%s", keyword)

	// 发送HTTP请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch baidu suggestions: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应（百度返回的是GBK编码）
	reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 解析返回数据
	// 格式: window.baidu.sug({q:"关键词",p:false,s:["建议1","建议2",...]});
	content := string(body)
	re := regexp.MustCompile(`s:\[([^\]]*)\]`)
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 || matches[1] == "" {
		return nil, nil // 没有建议词
	}

	// 解析关键词列表
	sugStr := matches[1]
	sugStr = strings.ReplaceAll(sugStr, `"`, "")
	suggestions := strings.Split(sugStr, ",")

	// 清理空白
	var result []string
	for _, s := range suggestions {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, s)
		}
	}

	return result, nil
}

// UpdateLangtailsIfNeeded 检查并更新长尾词（如果超过周期）
func UpdateLangtailsIfNeeded(sourceID int, sourceName string, cycleDays int) error {
	// 获取最新更新时间
	lastUptime, err := dao.GetLatestLangtailUptime(sourceID)
	if err != nil {
		// 如果查询失败，继续尝试抓取
		lastUptime = 0
	}

	// 计算周期（秒）
	cycleSeconds := int64(cycleDays * 24 * 3600)
	now := time.Now().Unix()

	// 如果未超过周期，直接返回
	if lastUptime > 0 && (now-lastUptime) < cycleSeconds {
		return nil
	}

	// 抓取新的长尾词
	suggestions, err := FetchBaiduSuggestions(sourceName)
	if err != nil {
		return err
	}

	// 插入数据库
	return dao.InsertLangtails(sourceID, sourceName, suggestions)
}
