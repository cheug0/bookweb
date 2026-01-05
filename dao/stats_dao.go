package dao

import "bookweb/utils"

// DashboardStats 仪表板统计数据
type DashboardStats struct {
	ArticleCount int
	UserCount    int
	TodayVisit   int
	SortCount    int
}

// GetDashboardStats 获取仪表板统计数据
func GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 小说总数
	utils.Db.QueryRow("SELECT COUNT(*) FROM jieqi_article_article").Scan(&stats.ArticleCount)

	// 用户总数
	utils.Db.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.UserCount)

	// 今日访问（所有小说的 dayvisit 汇总）
	utils.Db.QueryRow("SELECT COALESCE(SUM(dayvisit), 0) FROM jieqi_article_article").Scan(&stats.TodayVisit)

	// 分类数量
	utils.Db.QueryRow("SELECT COUNT(*) FROM sort").Scan(&stats.SortCount)

	return stats, nil
}
