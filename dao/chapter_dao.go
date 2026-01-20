package dao

import (
	"bookweb/model"
	"bookweb/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// GetChapterByID 根据ChapterID获取章节详情（包含内容）
func GetChapterByID(id int) (*model.Chapter, error) {
	var row *sql.Row
	if stmtGetChapterByID != nil {
		row = stmtGetChapterByID.QueryRow(id)
	} else {
		sqlStr := "select chapterid, siteid, articleid, articlename, volumeid, posterid, poster, postdate, lastupdate, chaptername, chapterorder, size, saleprice, salenum, totalcost, attachment, isvip, chaptertype, power, display from jieqi_article_chapter where chapterid = ?"
		row = utils.Db.QueryRow(sqlStr, id)
	}
	ch := &model.Chapter{}
	err := row.Scan(&ch.ChapterID, &ch.SiteID, &ch.ArticleID, &ch.ArticleName, &ch.VolumeID, &ch.PosterID, &ch.Poster, &ch.PostDate, &ch.LastUpdate, &ch.ChapterName, &ch.ChapterOrder, &ch.Size, &ch.SalePrice, &ch.SaleNum, &ch.TotalCost, &ch.Attachment, &ch.IsVIP, &ch.ChapterType, &ch.Power, &ch.Display)
	if err != nil {
		return nil, err
	}
	// Read content from file
	content, err := utils.GetChapterFileContent(ch.ArticleID, ch.ChapterID, ch.ChapterOrder)
	if err == nil {
		ch.Content = content
	} else {
		// 如果读取内容失败（文件不存在或路径错误），不应该返回 404，而是提示内容缺失
		ch.Content = "章节内容缺失，请联系管理员修复。"
		// fmt.Println("Error reading chapter content:", err) // Optional: Log to stdout/stderr
	}

	return ch, nil
}

// GetChapterByIDCached 带缓存获取章节详情（不包含内容文本）
// 缓存策略：仅缓存章节元数据，内容实时读取文件
func GetChapterByIDCached(id int) (*model.Chapter, error) {
	cacheKey := fmt.Sprintf("chapter_%d", id)

	var ch *model.Chapter

	// 尝试从缓存获取元数据
	if cached, err := utils.CacheGet(cacheKey); err == nil && cached != "" {
		ch = &model.Chapter{}
		if err := json.Unmarshal([]byte(cached), ch); err != nil {
			ch = nil // 解析失败，回源
		}
	}

	// 回源数据库
	if ch == nil {
		var err error
		ch, err = GetChapterByID(id)
		if err != nil {
			return nil, err
		}

		// 缓存元数据 (先清空 Content 防止大文本存入 Redis)
		realContent := ch.Content
		ch.Content = ""
		if data, err := json.Marshal(ch); err == nil {
			utils.CacheSet(cacheKey, string(data), 1*time.Hour)
		}
		ch.Content = realContent // 恢复内容以便返回
		return ch, nil
	}

	// 缓存命中后，需要读取文件内容
	content, err := utils.GetChapterFileContent(ch.ArticleID, ch.ChapterID, ch.ChapterOrder)
	if err == nil {
		ch.Content = content
	} else {
		ch.Content = "章节内容缺失，请联系管理员修复。"
	}

	return ch, nil
}

// GetPrevChapterID 获取上一章节ID
func GetPrevChapterID(articleID, currentOrder int) (int, error) {
	var id int
	var err error
	if stmtGetPrevChapterID != nil {
		err = stmtGetPrevChapterID.QueryRow(articleID, currentOrder).Scan(&id)
	} else {
		sqlStr := "select chapterid from jieqi_article_chapter where articleid = ? and chapterorder < ? order by chapterorder desc limit 1"
		err = utils.Db.QueryRow(sqlStr, articleID, currentOrder).Scan(&id)
	}
	return id, err
}

// GetPrevChapterIDCached 带缓存获取上一章节ID
func GetPrevChapterIDCached(articleID, currentOrder int) (int, error) {
	cacheKey := fmt.Sprintf("chapter_prev_%d_%d", articleID, currentOrder)
	if cached, err := utils.CacheGet(cacheKey); err == nil && cached != "" {
		if id, err := strconv.Atoi(cached); err == nil {
			return id, nil
		}
	}
	id, err := GetPrevChapterID(articleID, currentOrder)
	if err == nil {
		utils.CacheSet(cacheKey, strconv.Itoa(id), 1*time.Hour)
	}
	return id, err
}

// GetNextChapterIDCached 带缓存获取下一章节ID
func GetNextChapterIDCached(articleID, currentOrder int) (int, error) {
	cacheKey := fmt.Sprintf("chapter_next_%d_%d", articleID, currentOrder)
	if cached, err := utils.CacheGet(cacheKey); err == nil && cached != "" {
		if id, err := strconv.Atoi(cached); err == nil {
			return id, nil
		}
	}
	id, err := GetNextChapterID(articleID, currentOrder)
	if err == nil {
		utils.CacheSet(cacheKey, strconv.Itoa(id), 1*time.Hour)
	}
	return id, err
}

// GetNextChapterID 获取下一章节ID
func GetNextChapterID(articleID, currentOrder int) (int, error) {
	var id int
	var err error
	if stmtGetNextChapterID != nil {
		err = stmtGetNextChapterID.QueryRow(articleID, currentOrder).Scan(&id)
	} else {
		sqlStr := "select chapterid from jieqi_article_chapter where articleid = ? and chapterorder > ? order by chapterorder asc limit 1"
		err = utils.Db.QueryRow(sqlStr, articleID, currentOrder).Scan(&id)
	}
	return id, err
}
