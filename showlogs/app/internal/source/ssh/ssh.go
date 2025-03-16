package ssh

import (
	"bufio"
	"errors"
	"os/exec"
)

var errArgs = errors.New("command args invalid. ")

func LogString(args []string, callback func(string)) (err error) {
	if len(args) == 0 {
		return errArgs
	}
	cmd := exec.Command(args[0], args[1:]...)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer pipe.Close()

	// 错误输出重定向到stdout
	cmd.Stderr = cmd.Stdout

	if err = cmd.Start(); err != nil {
		return err
	}

	buf := bufio.NewScanner(pipe)
	for buf.Scan() {
		if callback != nil {
			callback(buf.Text())
		}
	}

	if err = cmd.Wait(); err != nil {
		return err
	}
	return nil
}
