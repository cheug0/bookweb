// bookmark.go
// 书签模型
// 定义用户阅读书签的数据结构
package model

// Bookmark 书签
type Bookmark struct {
	BookID       int
	ArticleID    int
	ArticleName  string
	UserID       int
	UserName     string
	ChapterID    int
	ChapterName  string
	ChapterOrder int
	JoinDate     int
}
