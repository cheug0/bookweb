// bookmark_service.go
// 书签服务
// 处理书签相关的业务逻辑
package service

import (
	"bookweb/dao"
	"bookweb/model" // 假设Utils里有时间获取
	"time"
)

// GetBookmarkList 获取用户的书签
func GetBookmarkList(userID, page, pageSize int) ([]*model.Bookmark, int, error) {
	offset := (page - 1) * pageSize
	bookmarks, err := dao.GetBookmarkList(userID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	count, err := dao.GetBookmarkCount(userID)
	if err != nil {
		return nil, 0, err
	}
	return bookmarks, count, nil
}

// AddToBookmark 添加书签
func AddToBookmark(userID, articleID, chapterID int) error {
	// 检查是否存在
	exists, err := dao.CheckBookmarkExists(userID, articleID, chapterID)
	if err != nil {
		return err
	}
	if exists {
		// 如果已存在，可以选择更新时间或者什么都不做
		// UpdateBookmarkTime(userID, articleID, chapterID) // 需要实现DAO
		return nil
	}

	// 获取文章信息
	article, err := dao.GetArticleByID(articleID)
	if err != nil {
		return err
	}
	if article == nil {
		return nil
	}

	// 获取章节信息（可做）

	user, err := dao.GetUserByID(userID)
	var userName string
	if err == nil && user != nil {
		userName = user.Username
	}

	// 构造书签对象
	bm := &model.Bookmark{
		ArticleID:   article.ArticleID,
		ArticleName: article.ArticleName,
		UserID:      userID,
		UserName:    userName,
		ChapterID:   chapterID,
		ChapterName: "", // 需要查询章节详情
		// ChapterOrder: 0,
		JoinDate: int(time.Now().Unix()),
	}

	// 补充ChapterName
	chapter, err := dao.GetChapterByID(chapterID)
	if err == nil && chapter != nil {
		bm.ChapterName = chapter.ChapterName
		bm.ChapterOrder = int(chapter.ChapterOrder)
	}

	return dao.AddBookmark(bm)
}

// RemoveFromBookmark 移出书签
func RemoveFromBookmark(bookID int) error {
	return dao.DeleteBookmark(bookID)
}
