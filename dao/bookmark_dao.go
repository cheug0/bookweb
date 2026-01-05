package dao

import (
	"bookweb/model"
	"bookweb/utils"
)

// GetBookmarkList 获取书签列表
func GetBookmarkList(userID int, offset, limit int) ([]*model.Bookmark, error) {
	sqlStr := "select bookid, articleid, articlename, userid, username, chapterid, chaptername, chapterorder, joindate from bookmark where userid = ? order by joindate desc limit ?, ?"
	rows, err := utils.Db.Query(sqlStr, userID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []*model.Bookmark
	for rows.Next() {
		bm := &model.Bookmark{}
		err := rows.Scan(&bm.BookID, &bm.ArticleID, &bm.ArticleName, &bm.UserID, &bm.UserName, &bm.ChapterID, &bm.ChapterName, &bm.ChapterOrder, &bm.JoinDate)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bm)
	}
	return bookmarks, nil
}

// GetBookmarkCount 获取书签总数
func GetBookmarkCount(userID int) (int, error) {
	sqlStr := "select count(*) from bookmark where userid = ?"
	var count int
	err := utils.Db.QueryRow(sqlStr, userID).Scan(&count)
	return count, err
}

// AddBookmark 添加书签
func AddBookmark(bm *model.Bookmark) error {
	sqlStr := "insert into bookmark(articleid, articlename, userid, username, chapterid, chaptername, chapterorder, joindate) values(?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := utils.Db.Exec(sqlStr, bm.ArticleID, bm.ArticleName, bm.UserID, bm.UserName, bm.ChapterID, bm.ChapterName, bm.ChapterOrder, bm.JoinDate)
	return err
}

// DeleteBookmark 删除书签
func DeleteBookmark(bookID int) error {
	sqlStr := "delete from bookmark where bookid = ?"
	_, err := utils.Db.Exec(sqlStr, bookID)
	return err
}

// CheckBookmarkExists 检查书签是否存在
func CheckBookmarkExists(userID, articleID, chapterID int) (bool, error) {
	sqlStr := "select count(*) from bookmark where userid = ? and articleid = ? and chapterid = ?"
	var count int
	err := utils.Db.QueryRow(sqlStr, userID, articleID, chapterID).Scan(&count)
	return count > 0, err
}

// DeleteBookmarkByArticle 删除指定章节书签
func DeleteBookmarkByChapter(userID, articleID, chapterID int) error {
	sqlStr := "delete from bookmark where userid = ? and articleid = ? and chapterid = ?"
	_, err := utils.Db.Exec(sqlStr, userID, articleID, chapterID)
	return err
}
