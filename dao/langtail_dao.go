package dao

import (
	"bookweb/model"
	"bookweb/utils"
	"time"
)

// GetLangtailsBySourceID 根据小说ID获取长尾词列表
func GetLangtailsBySourceID(sourceID int) ([]*model.Langtail, error) {
	sqlStr := "SELECT langid, sourceid, langname, sourcename, uptime FROM article_langtail WHERE sourceid = ? ORDER BY uptime DESC"
	rows, err := utils.Db.Query(sqlStr, sourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var langtails []*model.Langtail
	for rows.Next() {
		lt := &model.Langtail{}
		if err := rows.Scan(&lt.LangID, &lt.SourceID, &lt.LangName, &lt.SourceName, &lt.UpTime); err != nil {
			return nil, err
		}
		langtails = append(langtails, lt)
	}
	return langtails, nil
}

// GetLangtailByID 根据长尾词ID获取长尾词信息
func GetLangtailByID(langID int) (*model.Langtail, error) {
	sqlStr := "SELECT langid, sourceid, langname, sourcename, uptime FROM article_langtail WHERE langid = ?"
	row := utils.Db.QueryRow(sqlStr, langID)

	lt := &model.Langtail{}
	err := row.Scan(&lt.LangID, &lt.SourceID, &lt.LangName, &lt.SourceName, &lt.UpTime)
	if err != nil {
		return nil, err
	}
	return lt, nil
}

// InsertLangtails 批量插入长尾词（忽略重复）
func InsertLangtails(sourceID int, sourceName string, keywords []string) error {
	if len(keywords) == 0 {
		return nil
	}

	uptime := time.Now().Unix()
	sqlStr := "INSERT IGNORE INTO article_langtail (sourceid, langname, sourcename, uptime) VALUES (?, ?, ?, ?)"

	for _, kw := range keywords {
		if kw == "" {
			continue
		}
		_, err := utils.Db.Exec(sqlStr, sourceID, kw, sourceName, uptime)
		if err != nil {
			// 记录错误但继续执行
			continue
		}
	}
	return nil
}

// GetLatestLangtailUptime 获取指定小说的长尾词最新更新时间
func GetLatestLangtailUptime(sourceID int) (int64, error) {
	sqlStr := "SELECT MAX(uptime) FROM article_langtail WHERE sourceid = ?"
	var uptime *int64
	err := utils.Db.QueryRow(sqlStr, sourceID).Scan(&uptime)
	if err != nil {
		return 0, err
	}
	if uptime == nil {
		return 0, nil
	}
	return *uptime, nil
}
