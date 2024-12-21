package helper

import (
	"github.com/yzhlove/chat/chat05/global"
	"time"

	"github.com/jinzhu/now"
)

func toGameTime(t time.Time, offset int) time.Time {
	return t.In(time.FixedZone("", offset*3600))
}

// toGameDay returns the game day according to the configured refresh time and time zone offset
// t.Hour >= global.RefreshTime is Today and t.Hour < global.RefreshTime is Yesterday
func toGameDay(t time.Time, zone int) time.Time {
	return toGameTime(t, zone).Add(-global.RefreshTime * time.Hour)
}

// GameWeekDay consults zone and global.RefreshTime to get game Weekday
func GameWeekDay(t time.Time, zone int) time.Weekday {
	return toGameDay(t, zone).Weekday()
}

// GameYearDay consults zone and global.RefreshTime to get game YearDay
func GameYearDay(t time.Time, zone int) int {
	return toGameDay(t, zone).YearDay()
}

// GameWeek consults zone and global.RefreshTime to get game ISOWeek
func GameWeek(t time.Time, zone int) (year, week int) {
	return toGameDay(t, zone).ISOWeek()
}

// GameDate consults zone and global.RefreshTime to get game Date
func GameDate(t time.Time, zone int) (year int, month time.Month, day int) {
	return toGameDay(t, zone).Date()
}

func GameTime(t time.Time, zone int) time.Time {
	return toGameDay(t, zone)
}

// GameYearAndYearDay consults zone and global.RefreshTime to get game Year and YearDay
func GameYearAndYearDay(t time.Time, zone int) (int, int) {
	t = toGameDay(t, zone)
	return t.Year(), t.YearDay()
}

// IsSameGameDay 是否为同一刷新天
func IsSameGameDay(t1, t2 time.Time, timeZone int) bool {
	t1Year, _, t1Day := GameDate(t1, timeZone)
	t2Year, _, t2Day := GameDate(t2, timeZone)

	return t1Year == t2Year && t1Day == t2Day
}

// IsSameGameWeek 判断是否是同一周
func IsSameGameWeek(t1, t2 time.Time, timeZone int) bool {
	t1Year, t1Week := GameWeek(t1, timeZone)
	t2Year, t2Week := GameWeek(t2, timeZone)

	return t1Year == t2Year && t1Week == t2Week
}

// LastRefreshTime 获取rt类型，返回基于srvTime(UTC时间)的上次刷新时间点(UTC时间)
func LastRefreshTime(rt global.RefreshType, t time.Time, zone int) time.Time {
	return NextRefreshTime(rt, 0, t, zone)
}

// NextRefreshTime 获取rt类型，间隔interval，返回基于srvTime(UTC时间)的下一次刷新时间点(UTC时间)
func NextRefreshTime(rt global.RefreshType, round int, t time.Time, zone int) time.Time {
	t = toGameDay(t, zone)
	y, m, d := t.Date()

	switch rt {
	case global.Day:
		t = time.Date(y, m, d+round, global.RefreshTime, 0, 0, 0, t.Location())
	case global.Week:
		weekday := int(t.Weekday())
		if now.WeekStartDay != time.Sunday {
			weekStartDayInt := int(now.WeekStartDay)
			if weekday < weekStartDayInt {
				weekday = weekday + 7 - weekStartDayInt
			} else {
				weekday = weekday - weekStartDayInt
			}
		}
		t = time.Date(y, m, d-weekday+round*7, global.RefreshTime, 0, 0, 0, t.Location())
	case global.Month:
		t = time.Date(y, m+time.Month(round), 1, global.RefreshTime, 0, 0, 0, t.Location())
	}

	return t.UTC()
}

func BuildYearDay(year, yearDay int) uint32 {
	return uint32(year*1000 + yearDay)
}
