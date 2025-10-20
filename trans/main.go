package main

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"
)

func main() {

	var err error

	if err = InitConfig(); err != nil {
		slog.Error("init config failed! ", slog.Any("err", err))
		return
	}

	contents := make([]*Content, 0, 128)
	if err = ReadDir(func(c *Content) {
		c.optTime()
		contents = append(contents, c)
	}); err != nil {
		slog.Error("read dir failed! ", slog.Any("err", err))
		return
	}

	slices.SortFunc(contents, func(a, b *Content) int {
		return cmp.Compare(b.Timestamp, a.Timestamp)
	})

	for _, c := range contents {
		fmt.Println(c.String())
		fmt.Println()
	}
}
