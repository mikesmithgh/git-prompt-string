//go:build !windows

package integration

import "path/filepath"

var (
	notFoundMsg           string = "no such file or directory"
	git_prompt_string_bin string = "git-prompt-string"
)

func copyTestDataCmd(src string, dest string) (string, []string) {
	return "cp", []string{"-r", filepath.Join(src, "testdata"), dest}
}
