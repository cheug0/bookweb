package utils

import (
	"strings"
)

// IsBot checks if the User-Agent string belongs to a known bot/spider
func IsBot(userAgent string) bool {
	userAgent = strings.ToLower(userAgent)
	bots := []string{
		"bot",
		"spider",
		"crawl",
		"slurp",
		"mediapartners",
		"fast-webcrawler",
		"zyborg",
		"googlebot",
		"bingbot",
		"baiduspider",
		"yandexbot",
		"sogou",
		"exabot",
		"facebot",
		"ia_archiver",
	}

	for _, bot := range bots {
		if strings.Contains(userAgent, bot) {
			return true
		}
	}
	return false
}
