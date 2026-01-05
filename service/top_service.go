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

	allVisit, _ := dao.GetRankArticles("allvisit", limit)
	monthVisit, _ := dao.GetRankArticles("monthvisit", limit)
	weekVisit, _ := dao.GetRankArticles("weekvisit", limit)
	dayVisit, _ := dao.GetRankArticles("dayvisit", limit)

	allVote, _ := dao.GetRankArticles("allvote", limit)
	monthVote, _ := dao.GetRankArticles("monthvote", limit)
	weekVote, _ := dao.GetRankArticles("weekvote", limit)
	dayVote, _ := dao.GetRankArticles("dayvote", limit)

	newUpdate, _ := dao.GetRankArticles("lastupdate", limit)
	postDate, _ := dao.GetRankArticles("postdate", limit)

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
