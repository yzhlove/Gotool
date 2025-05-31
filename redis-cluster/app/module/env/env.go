package env

import (
	"os"
	"strings"
	"sync"
)

var mRedisPortStr string
var once sync.Once
var defRedisPorts = []string{"6381", "6382", "6383", "6384", "6385", "6386"}

func New() error {
	once.Do(func() {
		mRedisPortStr = os.Getenv("REDIS_PORTS")
	})
	return nil
}

func GetRedisPorts() []string {
	if len(mRedisPortStr) == 0 {
		return defRedisPorts
	}

	str := strings.Split(mRedisPortStr, ",")
	if len(str) != 6 {
		return defRedisPorts
	}
	return str
}
