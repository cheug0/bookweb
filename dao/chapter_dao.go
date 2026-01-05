package dao

import (
	"bookweb/model"
	"bookweb/utils"
)

// GetChapterByID 根据ChapterID获取章节详情（包含内容）
func GetChapterByID(id int) (*model.Chapter, error) {
	sqlStr := "select chapterid, siteid, articleid, articlename, volumeid, posterid, poster, postdate, lastupdate, chaptername, chapterorder, size, saleprice, salenum, totalcost, attachment, isvip, chaptertype, power, display from jieqi_article_chapter where chapterid = ?"
	row := utils.Db.QueryRow(sqlStr, id)
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

// GetPrevChapterID 获取上一章节ID
func GetPrevChapterID(articleID, currentOrder int) (int, error) {
	sqlStr := "select chapterid from jieqi_article_chapter where articleid = ? and chapterorder < ? order by chapterorder desc limit 1"
	var id int
	err := utils.Db.QueryRow(sqlStr, articleID, currentOrder).Scan(&id)
	return id, err
}

// GetNextChapterID 获取下一章节ID
func GetNextChapterID(articleID, currentOrder int) (int, error) {
	sqlStr := "select chapterid from jieqi_article_chapter where articleid = ? and chapterorder > ? order by chapterorder asc limit 1"
	var id int
	err := utils.Db.QueryRow(sqlStr, articleID, currentOrder).Scan(&id)
	return id, err
}
