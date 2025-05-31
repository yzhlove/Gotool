package tmpl

import (
	_ "embed"
	"fmt"
	"github.com/yzhlove/Gotool/redis-cluster/app/conf"
)

//go:embed redis-config.tpl
var RedisTpl string

type Redis struct {
	Port           string // redis 通信端口
	BusPort        string // redis 总线端口
	PidPath        string // pid 文件路径
	LogPath        string // log 文件路径
	DataDir        string // 持久化数据文件存放路径
	ClusterCfgName string // 集群配置文件名称
}

func NesRedisTpl(port string, functions ...RedisFunc) *Redis {
	redis := &Redis{
		Port:           port,
		BusPort:        fmt.Sprintf("1%s", port),
		PidPath:        fmt.Sprintf("/%s/%s/redis.pid", conf.RedisDir, port), // 默认路径 "/redis-cluster-test/port/redis.pid"
		LogPath:        fmt.Sprintf("/%s/%s/redis.log", conf.RedisDir, port), // 默认路径 "/redis-cluster-test/port/redis.log"
		DataDir:        fmt.Sprintf("/%s/%s/data", conf.RedisDir, port),      // 默认路径 "/redis-cluster-test/port/data"
		ClusterCfgName: fmt.Sprintf("nodes-%s.conf", port),                   // 默认名称 "nodes-port.conf"
	}
	for _, fn := range functions {
		fn(redis)
	}
	return redis
}

type RedisFunc func(tmpl *Redis)

func WithBusPort(busPort string) RedisFunc {
	return func(tmpl *Redis) {
		tmpl.BusPort = busPort
	}
}

func WithPIDPath(pidPath string) RedisFunc {
	return func(tmpl *Redis) {
		tmpl.PidPath = pidPath
	}
}

func WithLogPath(logPath string) RedisFunc {
	return func(tmpl *Redis) {
		tmpl.LogPath = logPath
	}
}

func WithDataDir(dataDir string) RedisFunc {
	return func(tmpl *Redis) {
		tmpl.DataDir = dataDir
	}
}

func WithClusterCfgName(clusterCfgName string) RedisFunc {
	return func(tmpl *Redis) {
		tmpl.ClusterCfgName = clusterCfgName
	}
}
