package model

// Sort 对应数据库中的 sort 表
type Sort struct {
	SortID    int    `json:"sortid"`    // 序号
	Weight    int    `json:"weight"`    // 排序
	Caption   string `json:"caption"`   // 分类名称
	ShortName string `json:"shortname"` // 分类简称
}
