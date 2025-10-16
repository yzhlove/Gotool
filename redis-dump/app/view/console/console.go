package console

import (
	"errors"
	"fmt"
	"os"

	"rain.com/Gotool/redis-dump/app/internal/rdb"
	"rain.com/Gotool/redis-dump/app/log"
)

func Console(path string) {

	if len(path) == 0 {
		log.Error("path is empty!")
		return
	}

	if _, err := os.Stat(path); err != nil {
		if !errors.As(err, &os.ErrNotExist) {
			log.Error("os.Stat failed!", log.ErrWrap(err))
			return
		}
	}

	f, err := os.Open(path)
	if err != nil {
		log.Error("os.Open failed!", log.ErrWrap(err))
		return
	}
	defer f.Close()

	metas, err := rdb.Dump(f)
	if err != nil {
		log.Error("parse file failed!", log.ErrWrap(err))
		return
	}

	for _, vv := range metas {
		fmt.Println(vv.String())
	}

}
