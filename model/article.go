package model

// Article 小说模型，对应 jieqi_article_article 表
type Article struct {
	ArticleID     int    `json:"articleid"`     // articleid
	SiteID        int    `json:"siteid"`        // siteid
	PostDate      int64  `json:"postdate"`      // postdate
	LastUpdate    int64  `json:"lastupdate"`    // lastupdate
	ArticleName   string `json:"articlename"`   // articlename
	Keywords      string `json:"keywords"`      // keywords
	Initial       string `json:"initial"`       // initial
	AuthorID      int    `json:"authorid"`      // authorid
	Author        string `json:"author"`        // author
	PosterID      int    `json:"posterid"`      // posterid
	Poster        string `json:"poster"`        // poster
	AgentID       int    `json:"agentid"`       // agentid
	Agent         string `json:"agent"`         // agent
	SortID        int    `json:"sortid"`        // sortid
	TypeID        int    `json:"typeid"`        // typeid
	Intro         string `json:"intro"`         // intro
	Notice        string `json:"notice"`        // notice
	Setting       string `json:"setting"`       // setting
	LastVolumeID  int    `json:"lastvolumeid"`  // lastvolumeid
	LastVolume    string `json:"lastvolume"`    // lastvolume
	LastChapterID int    `json:"lastchapterid"` // lastchapterid
	LastChapter   string `json:"lastchapter"`   // lastchapter
	Chapters      int    `json:"chapters"`      // chapters
	Size          int    `json:"size"`          // size
	LastVisit     int64  `json:"lastvisit"`     // lastvisit
	DayVisit      int    `json:"dayvisit"`      // dayvisit
	WeekVisit     int    `json:"weekvisit"`     // weekvisit
	MonthVisit    int    `json:"monthvisit"`    // monthvisit
	AllVisit      int    `json:"allvisit"`      // allvisit
	LastVote      int64  `json:"lastvote"`      // lastvote
	DayVote       int    `json:"dayvote"`       // dayvote
	WeekVote      int    `json:"weekvote"`      // weekvote
	MonthVote     int    `json:"monthvote"`     // monthvote
	AllVote       int    `json:"allvote"`       // allvote
	FullFlag      int    `json:"fullflag"`      // fullflag
	ImgFlag       int    `json:"imgflag"`       // imgflag
}
