package utils

import (
	"os"
	"path/filepath"
)

func IsDir(d string) (bool, error) {
	fi, err := os.Stat(d)
	if err != nil {
		return false, err
	}

	return fi.IsDir(), nil
}

func HasFileName(f string) bool {
	ext := filepath.Ext(f)
	return len(ext) > 0
}
