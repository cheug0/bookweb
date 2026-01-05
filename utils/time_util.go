package utils

import (
	"time"
)

// NowTime returns current unix timestamp
func NowTime() int64 {
	return time.Now().Unix()
}

// IsSameDay checks if two timestamps constitute the same day
func IsSameDay(t1, t2 int64) bool {
	tm1 := time.Unix(t1, 0)
	tm2 := time.Unix(t2, 0)
	return tm1.Year() == tm2.Year() && tm1.YearDay() == tm2.YearDay()
}

// IsSameWeek checks if two timestamps constitute the same week (ISO week)
func IsSameWeek(t1, t2 int64) bool {
	tm1 := time.Unix(t1, 0)
	tm2 := time.Unix(t2, 0)
	y1, w1 := tm1.ISOWeek()
	y2, w2 := tm2.ISOWeek()
	return y1 == y2 && w1 == w2
}

// IsSameMonth checks if two timestamps constitute the same month
func IsSameMonth(t1, t2 int64) bool {
	tm1 := time.Unix(t1, 0)
	tm2 := time.Unix(t2, 0)
	return tm1.Year() == tm2.Year() && tm1.Month() == tm2.Month()
}
