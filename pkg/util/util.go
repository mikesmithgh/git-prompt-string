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
	errorColor, _ := color.Color("red")
	clearColor, _ := color.Color("none")
	var error_msg string
	if e == nil {
		error_msg = "no error message provided"
	} else {
		error_msg = strings.ReplaceAll(e.Error(), "\n", "")
	}
	fmt.Printf("%s bgps error(%s): %s%s", errorColor, hint, error_msg, clearColor)
	os.Exit(exitCode)
}
