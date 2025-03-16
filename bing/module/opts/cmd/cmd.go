package cmd

import "flag"

type Values struct {
	UI bool
}

func New() *Values {
	values := &Values{}
	flag.BoolVar(&values.UI, "ui", true, "以图形界面的方式运行该工具!")

	flag.Parse()
	return values
}
