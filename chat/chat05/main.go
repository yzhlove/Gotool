package main

import (
	"fmt"
	"github.com/yzhlove/chat/chat05/global"
	"github.com/yzhlove/chat/chat05/signin"
	"os"
	"time"
)

func init() {
	os.Setenv("TZ", "UTC")
}

func main() {

	var cases = []struct {
		name  string
		cnt   time.Time
		recv  time.Time
		index int
		month int
		year  int
	}{
		{
			name:  "1",
			cnt:   buildTime("2024-12-01 04:59:59"),
			recv:  buildTime("2024-12-01 00:00:00"),
			index: 2,
			month: 11,
			year:  2024,
		},
		{
			name:  "2",
			cnt:   buildTime("2024-12-01 05:00:00"),
			recv:  buildTime("2024-12-01 00:00:00"),
			index: 1,
			month: 12,
			year:  2024,
		},
		{
			name:  "3",
			cnt:   buildTime("2024-12-01 05:00:00"),
			recv:  buildTime("2024-12-01 04:59:59"),
			index: 1,
			month: 12,
			year:  2024,
		},
		{
			name:  "4",
			cnt:   buildTime("2024-12-01 05:00:01"),
			recv:  buildTime("2024-12-01 00:00:00"),
			index: 1,
			month: 12,
			year:  2024,
		},
		{
			name:  "5",
			cnt:   buildTime("2024-12-01 05:00:01"),
			recv:  buildTime("2024-12-01 04:59:59"),
			index: 1,
			month: 12,
			year:  2024,
		},
		{
			name:  "6",
			cnt:   buildTime("2024-12-01 05:00:01"),
			recv:  buildTime("2024-12-01 05:00:00"),
			index: 2,
			month: 12,
			year:  2024,
		},
		{
			name:  "10",
			cnt:   buildTime("2024-12-01 05:00:05"),
			recv:  buildTime("2024-12-01 00:00:00"),
			index: 1,
			month: 12,
			year:  2024,
		},
		{
			name:  "11",
			cnt:   buildTime("2024-12-01 05:00:05"),
			recv:  buildTime("2024-12-01 04:59:59"),
			index: 1,
			month: 12,
			year:  2024,
		},
		{
			name:  "12",
			cnt:   buildTime("2024-12-01 05:00:05"),
			recv:  buildTime("2024-12-01 05:00:00"),
			index: 2,
			month: 12,
			year:  2024,
		},
		{
			name:  "13",
			cnt:   buildTime("2024-12-01 05:00:05"),
			recv:  buildTime("2024-12-01 05:00:01"),
			index: 2,
			month: 12,
			year:  2024,
		},
		{
			name:  "14",
			cnt:   buildTime("2024-12-31 04:59:59"),
			recv:  buildTime("2024-12-31 00:00:00"),
			index: 2,
			month: 12,
			year:  2024,
		},
		{
			name:  "15",
			cnt:   buildTime("2024-12-31 05:00:00"),
			recv:  buildTime("2024-12-31 04:59:59"),
			index: 2,
			month: 12,
			year:  2024,
		},
		{
			name:  "16",
			cnt:   buildTime("2024-12-31 05:00:01"),
			recv:  buildTime("2024-12-31 04:59:59"),
			index: 2,
			month: 12,
			year:  2024,
		},
		{
			name:  "17",
			cnt:   buildTime("2024-12-31 05:00:01"),
			recv:  buildTime("2024-12-31 05:00:00"),
			index: 2,
			month: 12,
			year:  2024,
		},
		{
			name:  "18",
			cnt:   buildTime("2025-01-01 04:59:59"),
			recv:  buildTime("2024-12-31 00:00:00"),
			index: 2,
			month: 12,
			year:  2024,
		},
		{
			name:  "19",
			cnt:   buildTime("2025-01-01 05:00:00"),
			recv:  buildTime("2024-12-31 04:59:59"),
			index: 1,
			month: 1,
			year:  2025,
		},
		{
			name:  "20",
			cnt:   buildTime("2025-01-01 05:00:01"),
			recv:  buildTime("2024-12-31 04:59:59"),
			index: 1,
			month: 1,
			year:  2025,
		},
		{
			name:  "21",
			cnt:   buildTime("2025-01-01 05:00:01"),
			recv:  buildTime("2024-12-31 05:00:00"),
			index: 1,
			month: 1,
			year:  2025,
		},
		{
			name:  "22",
			cnt:   buildTime("2025-01-01 05:00:00"),
			recv:  buildTime("2025-01-01 00:00:00"),
			index: 1,
			month: 1,
			year:  2025,
		},
		{
			name:  "23",
			cnt:   buildTime("2025-01-01 05:00:00"),
			recv:  buildTime("2025-01-01 04:59:59"),
			index: 1,
			month: 1,
			year:  2025,
		},
		{
			name:  "24",
			cnt:   buildTime("2025-01-01 05:00:01"),
			recv:  buildTime("2025-01-01 05:00:00"),
			index: 2,
			month: 1,
			year:  2025,
		},
	}

	for _, cas := range cases {
		y, m, i := signin.CalcTime2(cas.cnt, cas.recv, global.TimeZone)
		var status bool
		switch {
		case cas.year != y:
			status = true
			fallthrough
		case cas.month != m:
			status = true
			fallthrough
		case cas.index != i:
			status = true
		}
		if status {
			panic(fmt.Sprintf("{%s} test failed, [year:{%d,%d} month:{%d,%d} index:{%d,%d}] ", cas.name,
				cas.year, y, cas.month, m, cas.index, i))
		} else {
			fmt.Printf("{%s} is Ok.\n", cas.name)
		}
	}

}

func buildTime(timestr string) time.Time {

	loc, _ := time.LoadLocation("Asia/Shanghai")
	tm, _ := time.ParseInLocation(time.DateTime, timestr, loc)
	return tm
}

func test1() {
	now := time.Now()
	fmt.Println("now.String1 => ", now.String())
	fmt.Println("now.String1 => ", now.Format(time.RFC3339))
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	fmt.Println("now.String1 => ", now.In(loc).Format(time.RFC3339))
}

func test2() {

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}

	tm, err := time.ParseInLocation(time.DateTime, "2024-12-20 18:00:00", loc)
	if err != nil {
		panic(err)
	}

	fmt.Println("tm => ", tm.Format(time.RFC3339))
	fmt.Println("tm => ", tm.In(time.Local).Format(time.RFC3339))

}
