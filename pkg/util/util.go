package util

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/buildkite/shellwords"
	"github.com/mikesmithgh/git-prompt-string/pkg/color"
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
	return strings.TrimRight(string(result), "\r\n"), err
}

func ErrMsg(hint string, e error) {
	errorColor, _ := color.Color("red")
	clearColor, _ := color.Color("reset")
	var error_msg string
	if e == nil {
		error_msg = "no error message provided"
	} else {
		error_msg = strings.ReplaceAll(strings.ReplaceAll(e.Error(), "\n", ""), "\r", "")
	}
	fmt.Printf("%s git-prompt-string error(%s): %s%s", errorColor, hint, shellwords.Quote(error_msg), clearColor)
	os.Exit(1)
}
