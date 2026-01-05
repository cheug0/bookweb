package dao

import (
	"bookweb/model"
	"bookweb/utils"
)

// GetAllSorts 获取所有分类，按 weight 排序
func GetAllSorts() ([]*model.Sort, error) {
	sqlStr := "select sortid, weight, caption, shortname from sort order by weight asc"
	rows, err := utils.Db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sorts []*model.Sort
	for rows.Next() {
		s := &model.Sort{}
		err := rows.Scan(&s.SortID, &s.Weight, &s.Caption, &s.ShortName)
		if err != nil {
			return nil, err
		}
		sorts = append(sorts, s)
	}
	return sorts, nil
}

// GetSortByID 根据 ID 获取单个分类信息
func GetSortByID(sortID int) (*model.Sort, error) {
	sqlStr := "select sortid, weight, caption, shortname from sort where sortid = ?"
	row := utils.Db.QueryRow(sqlStr, sortID)
	s := &model.Sort{}
	err := row.Scan(&s.SortID, &s.Weight, &s.Caption, &s.ShortName)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// UpdateSort 更新分类信息
func UpdateSort(id int, caption, shortName string, weight int) error {
	sqlStr := "update sort set caption = ?, shortname = ?, weight = ? where sortid = ?"
	_, err := utils.Db.Exec(sqlStr, caption, shortName, weight, id)
	return err
}
