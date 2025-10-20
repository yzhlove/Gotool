package main

import "flag"

type Config struct {
	Path string
}

var _config *Config

func InitConfig() error {
	_config = new(Config)
	flag.StringVar(&_config.Path, "dir", "", "文件夹路径")
	flag.Parse()
	return nil
}

func GetConfig() *Config {
	return _config
}
