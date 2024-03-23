package integration

import "path/filepath"

var (
	notFoundMsg           string = "The system cannot find the path specified."
	git_prompt_string_bin string = "git-prompt-string.exe"
	escapedEqualSign      string = "^="
)

func copyTestDataCmd(src string, dest string) (string, []string) {
	return "xcopy", []string{"/S", "/E", "/I", filepath.Join(src, "testdata"), filepath.Join(dest, "testdata")}
}
