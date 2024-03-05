package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type GitRepo struct {
	GitDir                     string
	IsInGitDir                 bool
	IsInWorkTree               bool
	IsInBareRepo               bool
	IsInShallowRepo            bool
	IsSparseCheckout           bool
	Tag                        string
	AbbrevRef                  string
	ShortSha                   string
	PromptMergeStatus          string
	PromptSparseCheckoutStatus string
	PromptBranch               string
	PromptBareRepoStatus       string
}

func (g *GitRepo) GitDirFileExists(name string) (bool, error) {
	_, err := os.Stat(g.GitDirPath(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (g *GitRepo) IsGitDir(name string) bool {
	return IsDir(g.GitDirPath(name))
}

func (g *GitRepo) IsGitDirSymlink(name string) bool {
	return IsSymlink(g.GitDirPath(name))
}

func (g *GitRepo) GitDirPath(path string) string {
	return fmt.Sprintf("%s/%s", g.GitDir, path)
}

func (g *GitRepo) ReadGirDirFile(name string) (string, error) {
	return ReadFileTrimNewline(g.GitDirPath(name))
}

func (g *GitRepo) RevParseShort() ([]byte, error) {

	cmd := exec.Command(
		"git",
		"rev-parse",
		"--short",
		"@{upstream}",
	)

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	stderr, err := io.ReadAll(stderrPipe)
	if err != nil {
		return stderr, err
	}

	stdout, err := io.ReadAll(stdoutPipe)
	if err != nil {
		return stderr, err
	}

	err = cmd.Wait()

	if len(stdout) > 0 {
		g.ShortSha = strings.TrimSuffix(string(stdout), "\n")
	}

	return stderr, err
}

func (g *GitRepo) RevParse() ([]byte, error) {
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
		return nil, err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	stderr, err := io.ReadAll(stderrPipe)
	if err != nil {
		return stderr, err
	}

	stdout, err := io.ReadAll(stdoutPipe)
	if err != nil {
		return stderr, err
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
				shortStderr, shortErr := g.RevParseShort()
				err = errors.Join(err, shortErr)
				stderr = append(stderr, shortStderr...)
			}
		} else {
			return []byte{}, fmt.Errorf("expected result length of 5 or 6, got %d", resultLen)
		}
	}

	return stderr, err
}

type Result struct {
	stderr []byte
	err    error
}

func (g *GitRepo) BranchInfo() (string, error) {
	var err error
	ref := ""
	step := ""
	total := ""

	if g.IsGitDir("rebase-merge") {
		if ref, err = g.ReadGirDirFile("rebase-merge/head-name"); err != nil {
			return ref, err
		}
		if step, err = g.ReadGirDirFile("rebase-merge/msgnum"); err != nil {
			return ref, err
		}
		if total, err = g.ReadGirDirFile("rebase-merge/end"); err != nil {
			return ref, err
		}
		g.PromptMergeStatus = "|REBASE-m"
		if exists, exists_err := g.GitDirFileExists("rebase-merge/interactive"); exists {
			g.PromptMergeStatus = "|REBASE-i"
		} else if exists_err != nil {
			return strconv.FormatBool(exists), exists_err
		}

	} else {
		if g.IsGitDir("rebase-apply") {

			step, err = g.ReadGirDirFile("rebase-apply/next")
			if err != nil {
				return step, err
			}
			total, err = g.ReadGirDirFile("rebase-apply/last")
			if err != nil {
				return total, err
			}
			rebasing, err := g.GitDirFileExists("rebase-apply/rebasing")
			if err != nil {
				return strconv.FormatBool(rebasing), err
			}
			if rebasing {
				ref, err = g.ReadGirDirFile("rebase-apply/head-name")
				if err != nil {
					return ref, err
				}
				g.PromptMergeStatus = "|REBASE"
			} else {
				// TODO: check if we need to get branch name here, bgps was not doing it
				applying, err := g.GitDirFileExists("rebase-apply/applying")
				if err != nil {
					return strconv.FormatBool(applying), err
				}
				if applying {
					g.PromptMergeStatus = "|AM"
				} else {
					g.PromptMergeStatus = "|AM/REBASE"
				}
			}
		} else if g.IsGitDir("MERGE_HEAD") {
			g.PromptMergeStatus = "|MERGING"
		} else if g.IsGitDir("CHERRY_PICK_HEAD") {
			g.PromptMergeStatus = "|CHERRY-PICKING"
		} else if g.IsGitDir("REVERT_HEAD") {
			g.PromptMergeStatus = "|REVERTING"
		} else if g.IsGitDir("BISECT_LOG") {
			g.PromptMergeStatus = "|BISECTING"
		}

		if ref == "" {
			if g.IsGitDirSymlink("HEAD") {
				if ref, err = GitSymbolicRef("HEAD"); err != nil {
					return ref, err
				}
			} else {
				head := ""
				if head, err = g.ReadGirDirFile("HEAD"); err != nil {
					return head, err
				}
				ref = strings.TrimPrefix(head, "ref: ")

				if head == ref {
					tag, err := GitDescribeTag("HEAD")
					if err == nil {
						ref = tag
						g.Tag = ref
					} else if g.ShortSha == "" && len(head) > 7 {
						ref = head[:7]
					} else {
						ref = g.ShortSha
					}
					ref = fmt.Sprintf("(%s)", ref)
				}
			}
		}
	}

	if step != "" && total != "" {
		g.PromptMergeStatus += fmt.Sprintf(" %s/%s", step, total)
	}

	if g.PromptMergeStatus != "" {
		if _, err = GitLsFilesUnmerged(); err != nil {
			g.PromptMergeStatus += "|CONFLICT"
		}
	}

	if g.IsInGitDir {
		if g.IsInBareRepo {
			g.PromptBareRepoStatus = "BARE:"
		} else {
			ref = "GIT_DIR!"
		}
	}

	g.PromptBranch = strings.TrimPrefix(ref, "refs/heads/")

	// TODO: put behind a variable, default to false for perf
	checkSparse := false
	if checkSparse {
		g.IsSparseCheckout, err = GitSparseCheckout()
		if err != nil {
			errMsg("sparse checkout", err, 0)
		}

		if g.IsSparseCheckout {
			g.PromptSparseCheckoutStatus = "|SPARSE"
		}
	}

	prompt := fmt.Sprintf("%s%s%s%s", g.PromptBareRepoStatus, g.PromptBranch, g.PromptSparseCheckoutStatus, g.PromptMergeStatus)

	return prompt, nil
}

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

func GitSymbolicRef(ref string) (string, error) {
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

func GitDescribeTag(ref string) (string, error) {
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

func GitLsFilesUnmerged() (string, error) {
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

func GitSparseCheckout() (bool, error) {
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

func IsNoUpstreamErr(msg string) bool {
	amiguiousHead := "ambiguous argument 'HEAD'"
	noUpstream := "no upstream configured"
	noBranch := "no such branch"
	return strings.Contains(msg, amiguiousHead) || strings.Contains(msg, noUpstream) || strings.Contains(msg, noBranch)
}

func (g *GitRepo) GitHasCleanWorkingTree() (bool, error) {
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

func GitHasUntracked() (bool, error) {
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

func GitCommitCounts() (int, int, error) {
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

func errMsg(hint string, e error, exitCode int) {
	color.New(color.FgHiRed).Printf(" bgps error(%s): %s", hint, strings.ReplaceAll(e.Error(), "\n", ""))
	os.Exit(exitCode)
}

func main() {
	color.NoColor = false

	gitRepo := GitRepo{}

	stderr, err := gitRepo.RevParse()
	if err != nil {
		if strings.Contains(string(stderr), "not a git repository") {
			os.Exit(0)
		}
		// allow other errors to pass through, the git repo may not have upstream
	}

	branchInfo, err := gitRepo.BranchInfo()
	if err != nil {
		errMsg("branch info", err, 0)
	}

	prefix := "  "
	suffix := ""

	gitSymbol := ""

	if gitRepo.IsInBareRepo {
		gitSymbol = "󱣻"
	}

	if gitRepo.IsInBareRepo || gitRepo.IsInGitDir {
		c := color.New(color.FgHiBlack)
		if gitSymbol != "" {
			gitSymbol = " " + gitSymbol
		}
		c.Printf("%s%s%s%s", prefix, branchInfo, gitSymbol, suffix)
		return
	}

	cleanWorkingTree, err := gitRepo.GitHasCleanWorkingTree()
	if err != nil {
		errMsg("clean working tree", err, 0)
	}
	hasUntracked, err := GitHasUntracked()
	if err != nil {
		color.Red("%s", err)
		os.Exit(0)
	}

	ahead, behind := 0, 0
	if gitRepo.Tag == "" && gitRepo.ShortSha != "" {
		ahead, behind, err = GitCommitCounts()
	}
	if err != nil {
		color.New(color.FgRed).Printf("%s", err)
		os.Exit(0)
	}

	c := color.New()
	if cleanWorkingTree {
		c = color.New(color.FgGreen)
	}

	if ahead > 0 {
		c = color.New(color.FgYellow)
		gitSymbol = fmt.Sprintf("↑[%d]", ahead)
	}
	if behind > 0 {
		c = color.New(color.FgYellow)
		gitSymbol = fmt.Sprintf("↓[%d]", behind)
	}

	if ahead > 0 && behind > 0 {
		gitSymbol = fmt.Sprintf("↕ ↑[%d] ↓[%d]", ahead, behind)
	}

	if gitRepo.ShortSha == "" {
		c = color.New(color.FgHiBlack)
	}

	if hasUntracked {
		c = color.New(color.FgMagenta)
		gitSymbol = fmt.Sprintf("*%s", gitSymbol)
	}

	if !cleanWorkingTree && !hasUntracked {
		c = color.New(color.FgRed)
		gitSymbol = fmt.Sprintf("*%s", gitSymbol)
	}
	if gitSymbol != "" {
		gitSymbol = " " + gitSymbol
	}
	c.Printf("%s%s%s%s", prefix, branchInfo, gitSymbol, suffix)
}
