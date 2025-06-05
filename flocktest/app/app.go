package app

import (
	"fmt"
	"github.com/yzhlove/Gotool/flocktest/app/act"
	"github.com/yzhlove/Gotool/flocktest/app/cmd"
	"github.com/yzhlove/Gotool/flocktest/app/sh"
	"github.com/yzhlove/Gotool/flocktest/module/in"
	"github.com/yzhlove/Gotool/flocktest/module/log"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

func Run() error {

	input := in.Get()

	// 执行命令
	if len(input.Mode) != 0 {
		x := act.New(input.Path, input.Duration, input.Mode)
		return x.Act()
	}

	command := os.Args[0]
	lockPath := filepath.Join(filepath.Dir(command), "go.lock")

	// 执行指令
	c1 := cmd.Wrap("",
		cmd.Arg{Key: "-p", Vars: []string{lockPath}},
		cmd.Arg{Key: "-d", Vars: []string{"1"}},
		cmd.Arg{Key: "-m", Vars: []string{"1"}},
	)

	c2 := cmd.Wrap("",
		cmd.Arg{Key: "-p", Vars: []string{lockPath}},
		cmd.Arg{Key: "-d", Vars: []string{"7"}},
		cmd.Arg{Key: "-m", Vars: []string{"2"}},
	)

	c3 := cmd.Wrap("",
		cmd.Arg{Key: "-p", Vars: []string{lockPath}},
		cmd.Arg{Key: "-d", Vars: []string{"86"}},
		cmd.Arg{Key: "-m", Vars: []string{"3"}},
	)

	var wg sync.WaitGroup
	for k, args := range []cmd.Args{c1, c2, c3} {
		wg.Add(1)
		go func(k int, a cmd.Args) {
			defer wg.Done()
			if err := sh.Print(fmt.Sprintf("go.%d", k+1), command, a); err != nil {
				log.Error("command run error", slog.String("goroutine", fmt.Sprintf("go.%d", k+1)), log.ErrWrap(err))
			} else {
				log.Debug("command run ok. ", slog.String("goroutine", fmt.Sprintf("go.%d", k+1)))
			}
		}(k, args)
	}
	wg.Wait()
	return nil
}
