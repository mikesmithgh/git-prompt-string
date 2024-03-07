package util

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/mikesmithgh/bgps/pkg/color"
)

func IsDir(name string) bool {
	fileInfo, err := os.Stat(name)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func IsSymlink(name string) bool {
	fileInfo, err := os.Lstat(name)
	if err != nil {
		return false
	}
	return fileInfo.Mode()&fs.ModeSymlink != 0 // bitwise check if mode has symlink
}

func ReadFileTrimNewline(name string) (string, error) {
	result, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(result), "\n"), err
}

func ErrMsg(hint string, e error, exitCode int) {
	color, _ := color.Color("red")
	fmt.Printf("%s bgps error(%s): %s", color, hint, strings.ReplaceAll(e.Error(), "\n", ""))
	os.Exit(exitCode)
}
