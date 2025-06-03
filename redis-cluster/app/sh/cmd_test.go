package sh

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func Test_SH(t *testing.T) {

	fmt.Println(Which("redis-server"))

	fmt.Println(strings.HasSuffix("/usr/local/bin/redis-server", "redis-server"))
}

func Test_Node(t *testing.T) {

	cc := exec.Command("redis-server", "/Users/yostar/Desktop/MockRedis/redis-cluster-test/redis-meta/redis-6381.conf")
	res, err := cc.CombinedOutput()
	fmt.Println(err, string(res))

}
