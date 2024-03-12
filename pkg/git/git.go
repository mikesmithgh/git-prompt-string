package git

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

func CommitCounts() (int, int, error) {
	cmd := exec.Command(
		"git",
		"rev-list",
		"--left-right",
		"--count",
		"...@{upstream}",
	)
	stdCombined, err := cmd.CombinedOutput()
	if err != nil {
		return 0, 0, err
	}
	fields := strings.Fields(string(stdCombined))
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("expected field length of 2 got %d", len(fields))
	}
	ahead, _ := strconv.Atoi(fields[0])
	behind, _ := strconv.Atoi(fields[1])
	return ahead, behind, nil
}

func LsFilesUnmerged() (string, error) {
	cmd := exec.Command(
		"git",
		"ls-files",
		"--unmerged",
	)
	stdCombined, err := cmd.CombinedOutput()
	if err != nil {
		return string(stdCombined), err
	}
	return strings.TrimSuffix(string(stdCombined), "\n"), err
}

func SparseCheckout() (bool, error) {
	cmd := exec.Command(
		"git",
		"config",
		"--bool",
		"core.sparseCheckout",
	)
	stdCombined, err := cmd.CombinedOutput()
	if err != nil && len(stdCombined) != 0 {
		return false, err
	}
	isSparseCheckout, _ := strconv.ParseBool(strings.TrimSuffix(string(stdCombined), "\n"))
	return isSparseCheckout, nil
}

func SymbolicRef(ref string) (string, error) {
	cmd := exec.Command(
		"git",
		"symbolic-ref",
		ref,
	)
	stdCombined, err := cmd.CombinedOutput()
	if err != nil {
		return string(stdCombined), err
	}
	return strings.TrimSuffix(string(stdCombined), "\n"), err
}

func DescribeTag(ref string) (string, error) {
	cmd := exec.Command(
		"git",
		"describe",
		"--tags",
		"--exact-match",
		ref,
	)
	stdCombined, err := cmd.CombinedOutput()
	if err != nil {
		return string(stdCombined), err
	}
	return strings.TrimSuffix(string(stdCombined), "\n"), err
}

func IsNoUpstreamErr(msg string) bool {
	amiguiousHead := "ambiguous argument 'HEAD'"
	noUpstream := "no upstream configured"
	noBranch := "no such branch"
	return strings.Contains(msg, amiguiousHead) || strings.Contains(msg, noUpstream) || strings.Contains(msg, noBranch)
}

func HasUntracked() (bool, error) {
	exitCode := 0
	cmd := exec.Command(
		"git",
		"ls-files",
		"--others",
		"--exclude-standard",
		"--directory",
		"--no-empty-directory",
		"--error-unmatch",
		"--",
		":/*",
	)
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
	}
	if exitCode != 0 && exitCode != 1 {
		return false, err
	}
	return exitCode == 0, nil
}

func RevParseShort() (string, []byte, error) {

	cmd := exec.Command(
		"git",
		"rev-parse",
		"--short",
		"@{upstream}",
	)

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", nil, err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", nil, err
	}
	if err := cmd.Start(); err != nil {
		return "", nil, err
	}

	stderr, err := io.ReadAll(stderrPipe)
	if err != nil {
		return "", stderr, err
	}

	stdout, err := io.ReadAll(stdoutPipe)
	if err != nil {
		return "", stderr, err
	}

	err = cmd.Wait()

	return strings.TrimSuffix(string(stdout), "\n"), stderr, err
}

func RevParse() (*GitRepo, []byte, error) {
	g := GitRepo{}
	cmd := exec.Command(
		"git",
		"rev-parse",
		"--verify",
		"--absolute-git-dir",
		"--is-inside-git-dir",
		"--is-inside-work-tree",
		"--is-bare-repository",
		"--is-shallow-repository",
		"--abbrev-ref",
		"@{upstream}",
	)

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	stderr, err := io.ReadAll(stderrPipe)
	if err != nil {
		return nil, stderr, err
	}

	stdout, err := io.ReadAll(stdoutPipe)
	if err != nil {
		return nil, stderr, err
	}

	err = cmd.Wait()

	if len(stdout) > 0 {
		result := strings.Split(strings.TrimSuffix(string(stdout), "\n"), "\n")
		resultLen := len(result)
		if resultLen == 5 || resultLen == 6 {
			g.GitDir = result[0]
			g.IsInGitDir, _ = strconv.ParseBool(result[1])
			g.IsInWorkTree, _ = strconv.ParseBool(result[2])
			g.IsInBareRepo, _ = strconv.ParseBool(result[3])
			g.IsInShallowRepo, _ = strconv.ParseBool(result[4])
			if resultLen == 6 {
				g.AbbrevRef = result[5]
				shortSha, shortStderr, shortErr := RevParseShort()
				g.ShortSha = shortSha
				err = errors.Join(err, shortErr)
				stderr = append(stderr, shortStderr...)
			}
		} else {
			return nil, []byte{}, fmt.Errorf("expected result length of 5 or 6, got %d", resultLen)
		}
	}

	return &g, stderr, err
}

func HasCleanWorkingTree() (bool, error) {
	exitCode := 0
	cmd := exec.Command(
		"git",
		"diff",
		"--no-ext-diff",
		"--quiet",
		"HEAD",
	)
	stdCombined, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
		stderr := string(stdCombined)
		if IsNoUpstreamErr(stderr) {
			exitCode = 0
			// there is no upstream so compare against staging area
			cachedCmd := exec.Command(
				"git",
				"diff",
				"--cached",
				"--no-ext-diff",
				"--quiet",
			)
			err = cachedCmd.Run()
			if err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					exitCode = exitError.ExitCode()
				}
			}
		}
	}
	if exitCode != 0 && exitCode != 1 {
		return false, err
	}

	return exitCode == 0, nil
}

func BranchRemote(branch string) (string, error) {
	cmd := exec.Command(
		"git",
		"config",
		fmt.Sprintf("branch.%s.remote", branch),
	)
	stdCombined, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(stdCombined), "\n"), nil
}

func BranchMerge(branch string) (string, error) {
	cmd := exec.Command(
		"git",
		"config",
		fmt.Sprintf("branch.%s.merge", branch),
	)
	stdCombined, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(stdCombined), "\n"), nil
}
