package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var (
	builtBinaryPath string
	tmpDir          string
)

func TestMain(m *testing.M) {
	var err error
	tmpDir, err = os.MkdirTemp("", "integration")
	if err != nil {
		panic("failed to create temp dir")
	}
	defer os.RemoveAll(tmpDir)

	builtBinaryPath = filepath.Join(tmpDir, "bgps")

	cmd := exec.Command("go", "build", "-o", builtBinaryPath, "..")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("failed to build bgps: %s, %s", output, err))
	}

	cmd = exec.Command("cp", "-r", "../testdata", tmpDir)
	err = cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("failed to copy test data: %s", err))
	}

	test_dirs, err := os.ReadDir(filepath.Join(tmpDir, "testdata"))
	if err != nil {
		panic(fmt.Sprintf("failed to read test data dir: %s", err))
	}

	for _, test_dir := range test_dirs {
		test_dir_path := filepath.Join(tmpDir, "testdata", test_dir.Name())
		test_dir_files, err := os.ReadDir(test_dir_path)
		if err != nil {
			panic(fmt.Sprintf("failed to read test data file: %s", err))
		}
		for _, test_dir_file := range test_dir_files {
			if test_dir_file.Name() == "dot_git" {
				err = os.Rename(filepath.Join(test_dir_path, "dot_git"), filepath.Join(test_dir_path, ".git"))
				if err != nil {
					panic(fmt.Sprintf("failed to rename test data git dir: %s", err))
				}
			}
		}
	}
	fmt.Println("=== INIT")
	fmt.Println("tmpDir:", tmpDir)
	fmt.Println("builtBinaryPath:", builtBinaryPath)

	code := m.Run()
	defer os.Exit(code)
}
