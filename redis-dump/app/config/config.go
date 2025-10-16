package config

import "flag"

type Config struct {
	UI   bool
	Path string
}

func parse(conf *Config) {
	flag.BoolVar(&conf.UI, "ui", true, "已UI界面的形式启动!")
	flag.StringVar(&conf.Path, "p", "", "RedisDump文件路径!")
	flag.Parse()
}

func New() *Config {
	cc := &Config{}
	parse(cc)
	return cc
}
