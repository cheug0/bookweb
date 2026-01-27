// chapter.go
// 章节模型
// 定义小说章节内容及元数据结构
package model

// Chapter 章节模型，对应 jieqi_article_chapter 表
type Chapter struct {
	ChapterID    int    `json:"chapterid"`    // chapterid
	SiteID       int    `json:"siteid"`       // siteid
	ArticleID    int    `json:"articleid"`    // articleid
	ArticleName  string `json:"articlename"`  // articlename
	VolumeID     int    `json:"volumeid"`     // volumeid
	PosterID     int    `json:"posterid"`     // posterid
	Poster       string `json:"poster"`       // poster
	PostDate     int64  `json:"postdate"`     // postdate
	LastUpdate   int64  `json:"lastupdate"`   // lastupdate
	ChapterName  string `json:"chaptername"`  // chaptername
	ChapterOrder int    `json:"chapterorder"` // chapterorder
	Size         int    `json:"size"`         // size
	SalePrice    int    `json:"saleprice"`    // saleprice
	SaleNum      int    `json:"salenum"`      // salenum
	TotalCost    int    `json:"totalcost"`    // totalcost
	Attachment   string `json:"attachment"`   // attachment (content)
	IsVIP        int    `json:"isvip"`        // isvip
	ChapterType  int    `json:"chaptertype"`  // chaptertype
	Power        int    `json:"power"`        // power
	Display      int    `json:"display"`      // display
	Content      string `json:"content"`      // content from file
}
