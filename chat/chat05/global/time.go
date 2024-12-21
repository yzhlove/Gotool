package global

const (
	RefreshTime = 5 // 系统刷新时间点，早上五点
)

type RefreshType uint8

const (
	Day RefreshType = iota
	Week
	Month
)

const (
	TimeZone = 8
)
