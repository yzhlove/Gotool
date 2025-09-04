package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/yzhlove/Gotool/redis-cluster/app/conf"
	"github.com/yzhlove/Gotool/redis-cluster/app/helper"
	"github.com/yzhlove/Gotool/redis-cluster/app/module/env"
	"github.com/yzhlove/Gotool/redis-cluster/app/sh"
	"github.com/yzhlove/Gotool/redis-cluster/app/tmpl"
)

func Run() error {

	// 检查 redis 是否安装
	if err := whetherInstall(); err != nil {
		return err
	}

	workdir := env.GetWorkDir()
	// 模板文件存放路径
	var cfgDir string
	if len(workdir) != 0 {
		cfgDir = filepath.Join(workdir, conf.RedisDir, conf.RedisCfg)
	} else {
		cfgDir = filepath.Join(conf.RedisDir, conf.RedisCfg)
	}

	// 创建文件存放路径
	if err := helper.CreateDir(cfgDir); err != nil {
		return fmt.Errorf("create redis config dir error: %v", err)
	}

	// 创建 redis 模板
	t, err := template.New("redis").Parse(tmpl.RedisTemplate)
	if err != nil {
		return fmt.Errorf("parse redis template error: %v", err)
	}

	// 创建 redis 配置文件
	ports := env.GetRedisPorts()
	// 客户端IP
	clientIp := env.GetClientIp()

	for _, port := range ports {
		// 创建配置文件模板
		redisTpl := tmpl.NesRedisTpl(clientIp, port, tmpl.WithWorkDir(workdir))
		metapath := filepath.Join(cfgDir, fmt.Sprintf("redis-%s.conf", port))

		// 创建 redis 配置文件
		if err = writeTemplate(t, metapath, redisTpl); err != nil {
			return fmt.Errorf("write redis config file error: %v", err)
		}

		// 初始化 redis 目录环境
		if err = initRedisDir(redisTpl); err != nil {
			return fmt.Errorf("init redis dir error: %v", err)
		}

		// 根据配置文件启动 redis
		if err = sh.StartNode(metapath); err != nil {
			return fmt.Errorf("start redis node error: %v", err)
		}
	}

	// 加入 redis 集群
	if err = sh.StartCluster(clientIp, ports); err != nil {
		return fmt.Errorf("start redis cluster error: %v", err)
	}
	return nil
}

func whetherInstall() error {
	if sh.Which("redis-server") && sh.Which("redis-cli") {
		return nil
	}

	return fmt.Errorf("redis is not installed! ")
}

func writeTemplate(t *template.Template, path string, data any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, data)
}

func initRedisDir(t *tmpl.Redis) error {
	if err := helper.CreateDir(filepath.Dir(t.PidPath)); err != nil {
		return err
	}
	if err := helper.CreateDir(filepath.Dir(t.LogPath)); err != nil {
		return err
	}
	if err := helper.CreateDir(t.DataDir); err != nil {
		return err
	}
	return nil
}
