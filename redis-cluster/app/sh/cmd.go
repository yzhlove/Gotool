package sh

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/yzhlove/Gotool/redis-cluster/app/module/cmds"
	"github.com/yzhlove/Gotool/redis-cluster/app/module/log"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

func Which(command string) bool {
	cc := exec.Command("which", command)
	res, err := cc.Output()
	if err != nil {
		log.Warn("which command error", log.ErrWrap(err))
		return false
	}
	return strings.HasSuffix(string(res), command)
}

func StartNode(path string) error {
	cc := exec.Command("redis-server", path)
	res, err := cc.Output()
	if err != nil {
		log.Error("start node error", slog.String("path", path), log.ErrWrap(err))
		return err
	}
	if len(res) != 0 {
		err = errors.New(string(res))
		log.Error("start node error", slog.String("path", path), log.ErrWrap(err))
		return err
	}
	return nil
}

func StartCluster(ports []string) error {

	hosts := make([]string, 0, len(ports))
	for _, port := range ports {
		hosts = append(hosts, fmt.Sprintf("127.0.0.1:%s", port))
	}

	args := cmds.Wrap(
		"--cluster",
		cmds.Arg{
			Key: "create",
			Var: hosts,
		},
		cmds.Arg{
			Key: "--cluster-replicas",
			Var: []string{"1"},
		},
	)

	cc := exec.Command("redis-cli", args.V()...)

	input, err := cc.StdinPipe()
	if err != nil {
		log.Error("start pipeline cluster error", slog.Any("args", args.V()), log.ErrWrap(err))
		return err
	}
	defer input.Close()

	var output bytes.Buffer
	cc.Stdout = &output
	cc.Stdin = os.Stdin

	if err = cc.Start(); err != nil {
		log.Error("cmd.start cluster error", slog.Any("args", args.V()), log.ErrWrap(err))
		return err
	}

	// 输入参数
	io.WriteString(input, "yes\n")

	if err = cc.Wait(); err != nil {
		log.Error("cmd.wait cluster error", slog.Any("args", args.V()), log.ErrWrap(err))
		return err
	}

	// 执行结果
	result := output.Bytes()
	if !bytes.Contains(result, []byte("[OK] All 16384 slots covered.")) {
		log.Error("start cluster error",
			slog.Any("args", args.V()),
			slog.String("result", string(result)),
		)
		return fmt.Errorf("start cluster error")
	}
	return nil
}
