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
	SearchCacheTTL   = 3 * time.Minute
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

func searchCacheKey(keyword string, offset, limit int) string {
	return fmt.Sprintf("search:%s:%d:%d", keyword, offset, limit)
}

func searchCountCacheKey(keyword string) string {
	return fmt.Sprintf("search_count:%s", keyword)
}

// getCached 泛型缓存获取函数
func getCached[T any](key string, ttl time.Duration, fetchFunc func() (T, error)) (T, error) {
	// 1. Check Cache
	cached, err := utils.CacheGet(key)
	if err == nil && cached != "" {
		var result T
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result, nil
		}
	}

	// 2. Fetch
	result, err := fetchFunc()
	if err != nil {
		var zero T
		return zero, err
	}

	// 3. Set Cache
	if data, err := json.Marshal(result); err == nil {
		utils.CacheSet(key, string(data), ttl)
	}

	return result, nil
}

// GetArticleByIDCached 带缓存的获取文章
func GetArticleByIDCached(id int) (*model.Article, error) {
	if !utils.IsRedisEnabled() {
		return GetArticleByID(id)
	}

	return getCached(articleCacheKey(id), ArticleCacheTTL, func() (*model.Article, error) {
		return GetArticleByID(id)
	})
}

// GetChaptersByArticleIDCached 带缓存的获取章节列表
func GetChaptersByArticleIDCached(articleID int) ([]*model.Chapter, error) {
	if !utils.IsRedisEnabled() {
		return GetChaptersByArticleID(articleID)
	}

	return getCached(chaptersCacheKey(articleID), ChaptersCacheTTL, func() ([]*model.Chapter, error) {
		return GetChaptersByArticleID(articleID)
	})
}

// GetAllSortsCached 带缓存的获取所有分类
func GetAllSortsCached() ([]*model.Sort, error) {
	if !utils.IsRedisEnabled() {
		return GetAllSorts()
	}

	return getCached(sortsCacheKey(), SortsCacheTTL, func() ([]*model.Sort, error) {
		return GetAllSorts()
	})
}

// GetRankArticlesCached 带缓存的获取排行榜
func GetRankArticlesCached(orderBy string, limit int) ([]*model.Article, error) {
	if !utils.IsRedisEnabled() {
		return GetRankArticles(orderBy, limit)
	}

	return getCached(rankCacheKey(orderBy, limit), RankCacheTTL, func() ([]*model.Article, error) {
		return GetRankArticles(orderBy, limit)
	})
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

// SearchArticlesCached 带缓存的搜索文章
func SearchArticlesCached(keyword string, offset, limit int) ([]*model.Article, error) {
	if !utils.IsRedisEnabled() {
		return SearchArticles(keyword, offset, limit)
	}

	return getCached(searchCacheKey(keyword, offset, limit), SearchCacheTTL, func() ([]*model.Article, error) {
		return SearchArticles(keyword, offset, limit)
	})
}

// GetSearchCountCached 带缓存的获取搜索结果总数
func GetSearchCountCached(keyword string) (int, error) {
	if !utils.IsRedisEnabled() {
		return GetSearchCount(keyword)
	}

	return getCached(searchCountCacheKey(keyword), SearchCacheTTL, func() (int, error) {
		return GetSearchCount(keyword)
	})
}
