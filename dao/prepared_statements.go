package dao

import (
	"bookweb/utils"
	"database/sql"
)

// PreparedStatements 预编译的SQL语句
var (
	// Article 相关
	stmtGetArticleByID       *sql.Stmt
	stmtGetChaptersByArticle *sql.Stmt
	stmtGetVisitArticles     *sql.Stmt

	// Chapter 相关
	stmtGetChapterByID   *sql.Stmt
	stmtGetPrevChapterID *sql.Stmt
	stmtGetNextChapterID *sql.Stmt

	// Sort 相关
	stmtGetAllSorts *sql.Stmt
	stmtGetSortByID *sql.Stmt

	// User 相关
	stmtGetUserByUsername *sql.Stmt
	stmtGetUserByID       *sql.Stmt

	// Bookcase, Bookmark
	stmtGetBookcaseByUser *sql.Stmt
	stmtGetBookmarkByUser *sql.Stmt
)

// SQL语句常量
const (
	sqlGetArticleByID = `SELECT articleid, siteid, postdate, lastupdate, articlename, keywords, initial, authorid, author, posterid, poster, agentid, agent, sortid, typeid, intro, notice, setting, lastvolumeid, lastvolume, lastchapterid, lastchapter, chapters, size, lastvisit, dayvisit, weekvisit, monthvisit, allvisit, lastvote, dayvote, weekvote, monthvote, allvote, fullflag, imgflag FROM jieqi_article_article WHERE articleid = ?`

	sqlGetChaptersByArticle = `SELECT chapterid, chaptername, chapterorder, isvip, size, lastupdate FROM jieqi_article_chapter WHERE articleid = ? ORDER BY chapterorder ASC`

	sqlGetVisitArticles = `SELECT articleid, articlename, author, intro, size, lastupdate, sortid, fullflag, imgflag, lastchapterid, lastchapter FROM jieqi_article_article ORDER BY allvisit DESC LIMIT ?`

	sqlGetChapterByID = `SELECT chapterid, siteid, articleid, articlename, volumeid, posterid, poster, postdate, lastupdate, chaptername, chapterorder, size, saleprice, salenum, totalcost, attachment, isvip, chaptertype, power, display FROM jieqi_article_chapter WHERE chapterid = ?`

	sqlGetPrevChapterID = `SELECT chapterid FROM jieqi_article_chapter WHERE articleid = ? AND chapterorder < ? ORDER BY chapterorder DESC LIMIT 1`

	sqlGetNextChapterID = `SELECT chapterid FROM jieqi_article_chapter WHERE articleid = ? AND chapterorder > ? ORDER BY chapterorder ASC LIMIT 1`

	sqlGetAllSorts = `SELECT sortid, weight, caption, shortname FROM sort ORDER BY weight ASC`

	sqlGetSortByID = `SELECT sortid, weight, caption, shortname FROM sort WHERE sortid = ?`

	sqlGetUserByUsername = `SELECT id, username, password, email FROM users WHERE username = ?`

	sqlGetUserByID = `SELECT id, username, password, email, IFNULL(last_login_time, ''), IFNULL(current_login_time, '') FROM users WHERE id = ?`

	sqlGetBookcaseByUser = `SELECT caseid, articleid, articlename, userid, username, chapterid, chaptername, chapterorder, joindate, lastvisit, flag FROM bookcase WHERE userid = ? ORDER BY lastvisit DESC LIMIT ?, ?`

	sqlGetBookmarkByUser = `SELECT bookid, articleid, articlename, userid, username, chapterid, chaptername, chapterorder, joindate FROM bookmark WHERE userid = ? ORDER BY joindate DESC LIMIT ?, ?`
)

// InitPreparedStatements 初始化所有预编译语句
// 应在数据库连接建立后调用
func InitPreparedStatements() error {
	var err error
	db := utils.Db

	// Article 相关
	stmtGetArticleByID, err = db.Prepare(sqlGetArticleByID)
	if err != nil {
		return logPrepareError("GetArticleByID", err)
	}

	stmtGetChaptersByArticle, err = db.Prepare(sqlGetChaptersByArticle)
	if err != nil {
		return logPrepareError("GetChaptersByArticle", err)
	}

	stmtGetVisitArticles, err = db.Prepare(sqlGetVisitArticles)
	if err != nil {
		return logPrepareError("GetVisitArticles", err)
	}

	// Chapter 相关
	stmtGetChapterByID, err = db.Prepare(sqlGetChapterByID)
	if err != nil {
		return logPrepareError("GetChapterByID", err)
	}

	stmtGetPrevChapterID, err = db.Prepare(sqlGetPrevChapterID)
	if err != nil {
		return logPrepareError("GetPrevChapterID", err)
	}

	stmtGetNextChapterID, err = db.Prepare(sqlGetNextChapterID)
	if err != nil {
		return logPrepareError("GetNextChapterID", err)
	}

	// Sort 相关
	stmtGetAllSorts, err = db.Prepare(sqlGetAllSorts)
	if err != nil {
		return logPrepareError("GetAllSorts", err)
	}

	stmtGetSortByID, err = db.Prepare(sqlGetSortByID)
	if err != nil {
		return logPrepareError("GetSortByID", err)
	}

	// User 相关
	stmtGetUserByUsername, err = db.Prepare(sqlGetUserByUsername)
	if err != nil {
		return logPrepareError("GetUserByUsername", err)
	}

	stmtGetUserByID, err = db.Prepare(sqlGetUserByID)
	if err != nil {
		return logPrepareError("GetUserByID", err)
	}

	// Bookcase, Bookmark
	stmtGetBookcaseByUser, err = db.Prepare(sqlGetBookcaseByUser)
	if err != nil {
		return logPrepareError("GetBookcaseByUser", err)
	}

	stmtGetBookmarkByUser, err = db.Prepare(sqlGetBookmarkByUser)
	if err != nil {
		return logPrepareError("GetBookmarkByUser", err)
	}

	utils.LogInfo("DAO", "All prepared statements initialized successfully")
	return nil
}

func logPrepareError(name string, err error) error {
	utils.LogError("DAO", "Failed to prepare statement %s: %v", name, err)
	return err
}

// ClosePreparedStatements 关闭所有预编译语句
func ClosePreparedStatements() {
	stmts := []*sql.Stmt{
		stmtGetArticleByID, stmtGetChaptersByArticle, stmtGetVisitArticles,
		stmtGetChapterByID, stmtGetPrevChapterID, stmtGetNextChapterID,
		stmtGetAllSorts, stmtGetSortByID,
		stmtGetUserByUsername, stmtGetUserByID,
		stmtGetBookcaseByUser, stmtGetBookmarkByUser,
	}
	for _, stmt := range stmts {
		if stmt != nil {
			stmt.Close()
		}
	}
}
