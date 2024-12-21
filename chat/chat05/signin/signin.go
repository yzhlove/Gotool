package signin

import (
	"github.com/yzhlove/chat/chat05/global"
	"github.com/yzhlove/chat/chat05/helper"
	"time"
)

func CalcTime(now, receive time.Time, timezone int) (int, int, int) {

	var month = now.Month()
	var year = now.Year()
	var index = 1

	// 计算下一次的刷新时间
	if refreshTime := helper.NextRefreshTime(global.Month, 1, receive, timezone); now.Before(refreshTime) {
		forgetTime := refreshTime.AddDate(0, 0, -1)
		if forgetTime.Month() != refreshTime.Month() {
			if month == refreshTime.Month() {
				month = forgetTime.Month()
				year = forgetTime.Year()
			}
		}

		if forgetTime.Month() == refreshTime.Month() {
			if month != refreshTime.Month() {
				month = refreshTime.Month()
				year = refreshTime.Year()
			}
		}
		index++
	}
	return year, int(month), index
}

func CalcTime2(timeNow, receiveTime time.Time, timeZone int) (int, int, int) {

	// 计算当前时间的游戏时间
	gameTime := helper.GameTime(timeNow, timeZone)

	var month = gameTime.Month()
	var year = gameTime.Year()
	var index = 1

	refreshTime := helper.NextRefreshTime(global.Month, 1, receiveTime, timeZone)

	// 计算刷新时间的游戏时间
	refreshGameTime := helper.GameTime(refreshTime, timeZone)

	if gameTime.Before(refreshGameTime) {
		index++
	}
	return year, int(month), index
}
