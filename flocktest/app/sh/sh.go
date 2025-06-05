package sh

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/yzhlove/Gotool/flocktest/app/cmd"
	"io"
	"os/exec"
)

func Print(prefix, path string, args cmd.Args) error {

	cc := exec.Command(path, args.V()...)
	stdout, err := cc.StdoutPipe()
	if err != nil {
		return fmt.Errorf("[exec.StdoutPipe] error: %v", err)
	}

	stderr, err := cc.StderrPipe()
	if err != nil {
		return fmt.Errorf("[exec.StderrPipe] error: %v", err)
	}

	reader := bufio.NewReader(io.MultiReader(stdout, stderr))
	if err = cc.Start(); err != nil {
		return fmt.Errorf("[exec.Start] error: %v", err)
	}

	for {
		text, _, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("[reader.ReadLine] error: %v", err)
		}
		fmt.Printf("{%v} => %v \n", prefix, string(text))
	}

	if err = cc.Wait(); err != nil {
		return fmt.Errorf("[exec.Wait] error: %v", err)
	}
	return nil
}
