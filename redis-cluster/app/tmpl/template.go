package tmpl

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"github.com/yzhlove/Gotool/redis-cluster/app/conf"
)

//go:embed redis-config.tpl
var RedisTemplate string

type Redis struct {
	ClientIp       string // 指定的端口
	Port           string // redis 通信端口
	BusPort        string // redis 总线端口
	ClusterCfgName string // 集群配置文件名称
	WorkDir        string // 工作路径
	PidPath        string // pid 文件名称
	LogPath        string // 日志文件名称
	DataDir        string // 数据存储路径
}

func NesRedisTpl(clientIp, port string, functions ...RedisFunc) *Redis {
	redis := &Redis{
		ClientIp:       clientIp,
		Port:           port,
		BusPort:        fmt.Sprintf("1%s", port),
		PidPath:        fmt.Sprintf("/%s/%s/redis.pid", conf.RedisDir, port), // 默认路径 "/redis-cluster-test/port/redis.pid"
		LogPath:        fmt.Sprintf("/%s/%s/redis.log", conf.RedisDir, port), // 默认路径 "/redis-cluster-test/port/redis.log"
		DataDir:        fmt.Sprintf("/%s/%s/data", conf.RedisDir, port),      // 默认路径 "/redis-cluster-test/port/data"
		ClusterCfgName: fmt.Sprintf("nodes-%s.conf", port),                   // 默认名称 "nodes-port.conf"
		WorkDir:        "",
	}
	for _, fn := range functions {
		fn(redis)
	}
	return redis
}

type RedisFunc func(tmpl *Redis)

func WithWorkDir(dir string) RedisFunc {
	return func(tmpl *Redis) {
		tmpl.WorkDir = dir
		if len(dir) != 0 {
			tmpl.PidPath = filepath.Join(dir, tmpl.PidPath)
			tmpl.LogPath = filepath.Join(dir, tmpl.LogPath)
			tmpl.DataDir = filepath.Join(dir, tmpl.DataDir)
		}
	}
}
