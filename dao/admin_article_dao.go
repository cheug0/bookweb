package dao

import (
	"bookweb/model"
	"bookweb/utils"
	"fmt"
)

// GetArticleListAdmin 后台获取小说列表（分页+搜索）
func GetArticleListAdmin(page, pageSize int, keyword string) ([]*model.Article, int, error) {
	offset := (page - 1) * pageSize

	// 统计总数
	countSQL := "SELECT COUNT(*) FROM jieqi_article_article WHERE 1=1"
	args := []interface{}{}
	if keyword != "" {
		countSQL += " AND (articlename LIKE ? OR author LIKE ?)"
		args = append(args, "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int
	err := utils.Db.QueryRow(countSQL, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	querySQL := `SELECT articleid, articlename, author, sortid, intro, fullflag, 
		lastupdate, size, allvisit FROM jieqi_article_article WHERE 1=1`
	if keyword != "" {
		querySQL += " AND (articlename LIKE ? OR author LIKE ?)"
	}
	querySQL += " ORDER BY articleid DESC LIMIT ?, ?"
	args = append(args, offset, pageSize)

	rows, err := utils.Db.Query(querySQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var articles []*model.Article
	for rows.Next() {
		a := &model.Article{}
		err := rows.Scan(&a.ArticleID, &a.ArticleName, &a.Author, &a.SortID,
			&a.Intro, &a.FullFlag, &a.LastUpdate, &a.Size, &a.AllVisit)
		if err != nil {
			return nil, 0, err
		}
		articles = append(articles, a)
	}
	return articles, total, nil
}

// UpdateArticleAdmin 更新小说信息
func UpdateArticleAdmin(id int, name, author string, sortID, fullFlag int, intro string) error {
	sqlStr := `UPDATE jieqi_article_article SET 
		articlename = ?, author = ?, sortid = ?, fullflag = ?, intro = ? 
		WHERE articleid = ?`
	_, err := utils.Db.Exec(sqlStr, name, author, sortID, fullFlag, intro, id)
	return err
}

// DeleteArticleAdmin 删除小说
func DeleteArticleAdmin(id int) error {
	// 删除章节
	_, err := utils.Db.Exec("DELETE FROM jieqi_article_chapter WHERE articleid = ?", id)
	if err != nil {
		return fmt.Errorf("删除章节失败: %v", err)
	}
	// 删除小说
	_, err = utils.Db.Exec("DELETE FROM jieqi_article_article WHERE articleid = ?", id)
	if err != nil {
		return fmt.Errorf("删除小说失败: %v", err)
	}
	return nil
}

// GetArticleByIDAdmin 根据 ID 获取小说（后台用）
func GetArticleByIDAdmin(id int) (*model.Article, error) {
	sqlStr := `SELECT articleid, articlename, author, sortid, intro, fullflag 
		FROM jieqi_article_article WHERE articleid = ?`
	row := utils.Db.QueryRow(sqlStr, id)
	a := &model.Article{}
	err := row.Scan(&a.ArticleID, &a.ArticleName, &a.Author, &a.SortID, &a.Intro, &a.FullFlag)
	if err != nil {
		return nil, err
	}
	return a, nil
}
