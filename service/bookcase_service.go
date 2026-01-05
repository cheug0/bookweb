package service

import (
	"bookweb/dao"
	"bookweb/model"
	"time"
)

// GetBookcaseList 获取用户的书架
func GetBookcaseList(userID, page, pageSize int) ([]*model.Bookcase, int, error) {
	offset := (page - 1) * pageSize
	bookcases, err := dao.GetBookcaseList(userID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	count, err := dao.GetBookcaseCount(userID)
	if err != nil {
		return nil, 0, err
	}
	return bookcases, count, nil
}

// AddToBookcase 添加到书架
func AddToBookcase(userID, articleID int) error {
	// 检查是否存在
	exists, err := dao.CheckBookcaseExists(userID, articleID)
	if err != nil {
		return err
	}
	if exists {
		return nil // 已经存在，直接返回成功
	}

	// 获取文章信息
	article, err := dao.GetArticleByID(articleID)
	if err != nil {
		return err
	}
	if article == nil {
		return nil // 文章不存在
	}

	// 构造书架对象
	bc := &model.Bookcase{
		ArticleID:   article.ArticleID,
		ArticleName: article.ArticleName,
		UserID:      userID,
		// UserName:     "", // 可以从User表获取，或者Service层传入，这里暂时留空或不填
		ChapterID:   article.LastChapterID,
		ChapterName: article.LastChapter,
		// ChapterOrder: 0, // 可以查询Chapter表获取
		JoinDate:  int(time.Now().Unix()),
		LastVisit: int(time.Now().Unix()),
		Flag:      0,
	}
	// 补充UserName
	user, err := dao.GetUserByID(userID)
	if err == nil && user != nil {
		bc.UserName = user.Username
	}

	return dao.AddBookcase(bc)
}

// RemoveFromBookcase 移出书架
func RemoveFromBookcase(userID, articleID int) error {
	return dao.DeleteBookcaseByArticle(userID, articleID)
}

// RemoveBookcaseByID 根据ID移出
func RemoveBookcaseByID(caseID int) error {
	return dao.DeleteBookcase(caseID)
}
