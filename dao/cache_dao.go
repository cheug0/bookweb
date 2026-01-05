package dao

import (
	"bookweb/model"
	"bookweb/utils"
	"encoding/json"
	"fmt"
	"time"
)

// Cache TTL constants
const (
	ArticleCacheTTL  = 5 * time.Minute
	ChaptersCacheTTL = 10 * time.Minute
	SortsCacheTTL    = 30 * time.Minute
	RankCacheTTL     = 5 * time.Minute
)

// Cache key generators
func articleCacheKey(id int) string {
	return fmt.Sprintf("article:%d", id)
}

func chaptersCacheKey(articleID int) string {
	return fmt.Sprintf("chapters:%d", articleID)
}

func sortsCacheKey() string {
	return "sorts:all"
}

func rankCacheKey(orderBy string, limit int) string {
	return fmt.Sprintf("rank:%s:%d", orderBy, limit)
}

// GetArticleByIDCached 带缓存的获取文章
func GetArticleByIDCached(id int) (*model.Article, error) {
	if !utils.IsRedisEnabled() {
		return GetArticleByID(id)
	}

	key := articleCacheKey(id)
	cached, err := utils.CacheGet(key)
	if err == nil && cached != "" {
		var article model.Article
		if json.Unmarshal([]byte(cached), &article) == nil {
			return &article, nil
		}
	}

	// Cache miss, get from DB
	article, err := GetArticleByID(id)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if data, err := json.Marshal(article); err == nil {
		utils.CacheSet(key, string(data), ArticleCacheTTL)
	}

	return article, nil
}

// GetChaptersByArticleIDCached 带缓存的获取章节列表
func GetChaptersByArticleIDCached(articleID int) ([]*model.Chapter, error) {
	if !utils.IsRedisEnabled() {
		return GetChaptersByArticleID(articleID)
	}

	key := chaptersCacheKey(articleID)
	cached, err := utils.CacheGet(key)
	if err == nil && cached != "" {
		var chapters []*model.Chapter
		if json.Unmarshal([]byte(cached), &chapters) == nil {
			return chapters, nil
		}
	}

	// Cache miss
	chapters, err := GetChaptersByArticleID(articleID)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(chapters); err == nil {
		utils.CacheSet(key, string(data), ChaptersCacheTTL)
	}

	return chapters, nil
}

// GetAllSortsCached 带缓存的获取所有分类
func GetAllSortsCached() ([]*model.Sort, error) {
	if !utils.IsRedisEnabled() {
		return GetAllSorts()
	}

	key := sortsCacheKey()
	cached, err := utils.CacheGet(key)
	if err == nil && cached != "" {
		var sorts []*model.Sort
		if json.Unmarshal([]byte(cached), &sorts) == nil {
			return sorts, nil
		}
	}

	sorts, err := GetAllSorts()
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(sorts); err == nil {
		utils.CacheSet(key, string(data), SortsCacheTTL)
	}

	return sorts, nil
}

// GetRankArticlesCached 带缓存的获取排行榜
func GetRankArticlesCached(orderBy string, limit int) ([]*model.Article, error) {
	if !utils.IsRedisEnabled() {
		return GetRankArticles(orderBy, limit)
	}

	key := rankCacheKey(orderBy, limit)
	cached, err := utils.CacheGet(key)
	if err == nil && cached != "" {
		var articles []*model.Article
		if json.Unmarshal([]byte(cached), &articles) == nil {
			return articles, nil
		}
	}

	articles, err := GetRankArticles(orderBy, limit)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(articles); err == nil {
		utils.CacheSet(key, string(data), RankCacheTTL)
	}

	return articles, nil
}

// InvalidateArticleCache 使文章缓存失效
func InvalidateArticleCache(id int) {
	if utils.IsRedisEnabled() {
		utils.CacheDel(articleCacheKey(id))
		utils.CacheDel(chaptersCacheKey(id))
	}
}

// InvalidateSortsCache 使分类缓存失效
func InvalidateSortsCache() {
	if utils.IsRedisEnabled() {
		utils.CacheDel(sortsCacheKey())
	}
}
