package model

// Langtail 长尾词数据模型
type Langtail struct {
	LangID     int    `json:"langid"`
	SourceID   int    `json:"sourceid"`
	LangName   string `json:"langname"`
	SourceName string `json:"sourcename"`
	UpTime     int64  `json:"uptime"`
}
