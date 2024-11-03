package db

import "github.com/yzhlove/upload/app/entity"

type DBer interface {
	Get(filename string) (entity.FileMeta, error)
}
