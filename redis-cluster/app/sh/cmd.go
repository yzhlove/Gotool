package sh

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	"github.com/yzhlove/Gotool/redis-cluster/app/module/cmds"
	"github.com/yzhlove/Gotool/redis-cluster/app/module/log"
)

func Which(command string) bool {
	_, err := exec.LookPath(command)
	if err != nil {
		log.Warn("which command error", log.ErrWrap(err))
		return false
	}
	return true
}

func StartNode(path string) error {
	cc := exec.Command("redis-server", path)
	res, err := cc.CombinedOutput()
	if err != nil {
		log.Error("start node error",
			slog.String("path", path),
			log.ErrWrap(err),
			slog.String("reason", string(res)),
		)
		return err
	}
	return nil
}

func StartCluster(clientIp string, ports []string) error {

	hosts := make([]string, 0, len(ports))
	for _, port := range ports {
		hosts = append(hosts, fmt.Sprintf("%s:%s", clientIp, port))
	}

	args := cmds.NewCMD("",
		cmds.Arg{
			Key: "--cluster",
		},
		cmds.Arg{
			Key: "create",
			Var: hosts,
		},
		cmds.Arg{
			Key: "--cluster-replicas",
			Var: []string{"1"},
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	cc := exec.CommandContext(ctx, "redis-cli", args.V()...)
	input, err := cc.StdinPipe()
	if err != nil {
		log.Error("get stdin pipe error", slog.Any("args", args.V()), log.ErrWrap(err))
		return err
	}

	output, err := cc.StdoutPipe()
	if err != nil {
		log.Error("get stdout pipe error", slog.Any("args", args.V()), log.ErrWrap(err))
		return err
	}

	errput, err := cc.StderrPipe()
	if err != nil {
		log.Error("get stderr pipe error", slog.Any("args", args.V()), log.ErrWrap(err))
		return err
	}

	var buffer = bufio.NewReader(io.MultiReader(output, errput))
	var store bytes.Buffer

	// 打印执行解锁
	defer func() {
		fmt.Println("command: redis-cli ", strings.Join(args.V(), " "))
		fmt.Println(store.String())
	}()

	if err = cc.Start(); err != nil {
		log.Error("cmd.start cluster error", slog.Any("args", args.V()), log.ErrWrap(err))
		return err
	}

	// 输入参数
	go func() {
		defer input.Close()
		time.Sleep(time.Millisecond * 50)
		io.WriteString(input, "yes\n")
	}()

	var cache = make([]byte, 512)
LOOP:
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("cmd.start failed: context with cancel! ")
		default:
			n, err := buffer.Read(cache)
			if err != nil {
				if errors.Is(err, io.EOF) {
					store.Write(cache[:n])
					break LOOP
				}
				return fmt.Errorf("cmd.satart failed: %v", err)
			}
			store.Write(cache[:n])
		}
	}

	if err = cc.Wait(); err != nil {
		log.Error("cmd.wait cluster error", slog.Any("args", args.V()), log.ErrWrap(err))
		return err
	}

	// 执行结果
	result := store.Bytes()
	if !bytes.Contains(result, []byte("[OK] All 16384 slots covered.")) {
		log.Error("start cluster error",
			slog.Any("args", args.V()),
			slog.String("result", string(result)),
		)
		return fmt.Errorf("start redis cluster error")
	}
	return nil
}
