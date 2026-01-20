package dao

import (
	"bookweb/model"
	"bookweb/utils"
	"database/sql"
)

// GetBookcaseList 获取书架列表
func GetBookcaseList(userID int, offset, limit int) ([]*model.Bookcase, error) {
	var rows *sql.Rows
	var err error
	if stmtGetBookcaseByUser != nil {
		rows, err = stmtGetBookcaseByUser.Query(userID, offset, limit)
	} else {
		sqlStr := "select caseid, articleid, articlename, userid, username, chapterid, chaptername, chapterorder, joindate, lastvisit, flag from bookcase where userid = ? order by lastvisit desc limit ?, ?"
		rows, err = utils.Db.Query(sqlStr, userID, offset, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookcases []*model.Bookcase
	for rows.Next() {
		bc := &model.Bookcase{}
		err := rows.Scan(&bc.CaseID, &bc.ArticleID, &bc.ArticleName, &bc.UserID, &bc.UserName, &bc.ChapterID, &bc.ChapterName, &bc.ChapterOrder, &bc.JoinDate, &bc.LastVisit, &bc.Flag)
		if err != nil {
			return nil, err
		}
		bookcases = append(bookcases, bc)
	}
	return bookcases, nil
}

// GetBookcaseCount 获取书架总数
func GetBookcaseCount(userID int) (int, error) {
	sqlStr := "select count(*) from bookcase where userid = ?"
	var count int
	err := utils.Db.QueryRow(sqlStr, userID).Scan(&count)
	return count, err
}

// AddBookcase 添加到书架
func AddBookcase(bc *model.Bookcase) error {
	sqlStr := "insert into bookcase(articleid, articlename, userid, username, chapterid, chaptername, chapterorder, joindate, lastvisit, flag) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := utils.Db.Exec(sqlStr, bc.ArticleID, bc.ArticleName, bc.UserID, bc.UserName, bc.ChapterID, bc.ChapterName, bc.ChapterOrder, bc.JoinDate, bc.LastVisit, bc.Flag)
	return err
}

// DeleteBookcase 删除书架
func DeleteBookcase(caseID int) error {
	sqlStr := "delete from bookcase where caseid = ?"
	_, err := utils.Db.Exec(sqlStr, caseID)
	return err
}

// CheckBookcaseExists 检查是否存在
func CheckBookcaseExists(userID, articleID int) (bool, error) {
	sqlStr := "select count(*) from bookcase where userid = ? and articleid = ?"
	var count int
	err := utils.Db.QueryRow(sqlStr, userID, articleID).Scan(&count)
	return count > 0, err
}

// DeleteBookcaseByArticle 删除指定类型
func DeleteBookcaseByArticle(userID, articleID int) error {
	sqlStr := "delete from bookcase where userid = ? and articleid = ?"
	_, err := utils.Db.Exec(sqlStr, userID, articleID)
	return err
}
