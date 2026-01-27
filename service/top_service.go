// top_service.go
// 排行榜服务
// 处理排行榜数据的聚合与计算
package service

import (
	"bookweb/dao"
	"bookweb/model"
)

// TopData 包含排行榜页面的所有数据
type TopData struct {
	AllVisit   []*model.Article
	MonthVisit []*model.Article
	WeekVisit  []*model.Article
	DayVisit   []*model.Article
	AllVote    []*model.Article
	MonthVote  []*model.Article
	WeekVote   []*model.Article
	DayVote    []*model.Article
	NewUpdate  []*model.Article
	PostDate   []*model.Article
}

// GetTopData 获取排行榜汇总数据
func GetTopData() (*TopData, error) {
	limit := 10 // 每类取前10名

	allVisit, _ := dao.GetRankArticlesCached("allvisit", limit)
	monthVisit, _ := dao.GetRankArticlesCached("monthvisit", limit)
	weekVisit, _ := dao.GetRankArticlesCached("weekvisit", limit)
	dayVisit, _ := dao.GetRankArticlesCached("dayvisit", limit)

	allVote, _ := dao.GetRankArticlesCached("allvote", limit)
	monthVote, _ := dao.GetRankArticlesCached("monthvote", limit)
	weekVote, _ := dao.GetRankArticlesCached("weekvote", limit)
	dayVote, _ := dao.GetRankArticlesCached("dayvote", limit)

	newUpdate, _ := dao.GetRankArticlesCached("lastupdate", limit)
	postDate, _ := dao.GetRankArticlesCached("postdate", limit)

	return &TopData{
		AllVisit:   allVisit,
		MonthVisit: monthVisit,
		WeekVisit:  weekVisit,
		DayVisit:   dayVisit,
		AllVote:    allVote,
		MonthVote:  monthVote,
		WeekVote:   weekVote,
		DayVote:    dayVote,
		NewUpdate:  newUpdate,
		PostDate:   postDate,
	}, nil
}
