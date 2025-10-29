package main

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func ReadDir(callback func(c *Content)) error {

	cc := GetConfig()
	if len(cc.Path) == 0 {
		slog.Warn("path is empty! ")
		return nil
	}

	if !checkDir(cc.Path) {
		slog.Warn("path must be is directory! ")
		return nil
	}

	chstr := make(chan string, 16)
	stop := make(chan struct{})
	go func() {
		defer close(stop)
		for content := range chstr {
			if c := toContent(content); c != nil {
				c.optTime()
				callback(c)
			}
		}
	}()

	var wg sync.WaitGroup
	if err := filepath.WalkDir(cc.Path, func(path string, d fs.DirEntry, err error) error {
		if err == nil {
			if !d.IsDir() {
				if strings.HasPrefix(d.Name(), ".") || strings.HasPrefix(d.Name(), "~") {
					return nil
				}
				wg.Go(func() {
					value, err := toJson(path)
					if err != nil {
						slog.Error("to json failed! ", slog.String("path", path), slog.Any("error", err))
						return
					}
					if strs := getContent(value); len(strs) != 0 {
						for _, meta := range strs {
							if len(meta) > 0 {
								chstr <- meta
							}
						}
					}
				})
			}
		}
		return nil
	}); err != nil {
		return err
	}
	wg.Wait()
	close(chstr)
	<-stop
	return nil
}

func checkDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		slog.Warn("os.Stat failed ", slog.Any("error", err))
		return false
	}
	return s.IsDir()
}
