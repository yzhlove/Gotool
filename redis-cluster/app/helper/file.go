package helper

import "os"

func CreateDir(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
