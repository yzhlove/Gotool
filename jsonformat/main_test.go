package main

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"
	"testing"
)

func Test_Read(t *testing.T) {

	_config = new(Config)
	_config.Path = ""

	var err error
	var contents = make([]*Content, 0, 128)
	if err = ReadDir(func(c *Content) {
		contents = append(contents, c)
	}); err != nil {
		slog.Error("read dir failed! ", slog.Any("err", err))
		return
	}

	//slices.SortFunc(contents, func(a, b *Content) int {
	//	return cmp.Compare(b.Timestamp, a.Timestamp)
	//})

	//fmt.Println("log count: ", len(contents))
	//
	//for _, x := range contents {
	//	fmt.Println(x.String())
	//}

	showUserPath(contents)

}

func showUserPath(contents []*Content) {
	// 升序排列
	slices.SortFunc(contents, func(a, b *Content) int { return cmp.Compare(a.Timestamp, b.Timestamp) })

	meta := new(userActionMeta)

	for _, cc := range contents {
		meta.push(&userAction{
			time:    cc.Time,
			req:     cc.Req,
			msg:     cc.Msg,
			isError: cc.Level == "ERROR",
		})
	}

	for _, cc := range meta.queue {
		if cc.isError {
			fmt.Printf("xxxxx>time:%s ERROR req:%s msg:%s \n", cc.time, cc.req, cc.msg)
		} else {
			fmt.Printf("=====>time:%s req:%s count:%d msg:%s \n", cc.time, cc.req, cc.count, cc.msg)
		}
	}
}

type userAction struct {
	time    string
	req     string
	msg     string
	count   int
	isError bool
}

type userActionMeta struct {
	queue []*userAction
}

func (meta *userActionMeta) push(action *userAction) {

	if action.isError || len(meta.queue) == 0 {
		meta.queue = append(meta.queue, action)
	} else {
		end := meta.queue[len(meta.queue)-1]
		if end.req == action.req {
			end.count++
		} else {
			meta.queue = append(meta.queue, action)
		}
	}
}
