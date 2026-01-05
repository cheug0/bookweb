package model

// Bookcase 书架
type Bookcase struct {
	CaseID       int
	ArticleID    int
	ArticleName  string
	UserID       int
	UserName     string
	ChapterID    int
	ChapterName  string
	ChapterOrder int
	JoinDate     int
	LastVisit    int
	Flag         int
}
