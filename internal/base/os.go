// Copyright 2023 enpsl. All rights reserved.

// file path op func

package base

import (
	"github.com/enpsl/conf-reload/internal/errors"
	"os"
	"strings"
)

func FindParentDir(path string) (string, error) {
	isDir, err := isDirectory(path)
	if err != nil || isDir {
		return path, errors.ErrFormat(errors.ErrInvalidFilePath, err)
	}
	return getParentDirectory(path), nil
}

func isDirectory(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	switch mode := f.Mode(); {
	case mode.IsDir():
		return true, nil
	case mode.IsRegular():
		return false, nil
	}
	return false, nil
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
