package git

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mikesmithgh/bgps/pkg/color"
	"github.com/mikesmithgh/bgps/pkg/config"
	"github.com/mikesmithgh/bgps/pkg/util"
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
	return util.IsDir(g.GitDirPath(name))
}

func (g *GitRepo) IsGitDirSymlink(name string) bool {
	return util.IsSymlink(g.GitDirPath(name))
}

func (g *GitRepo) GitDirPath(path string) string {
	return fmt.Sprintf("%s/%s", g.GitDir, path)
}

func (g *GitRepo) ReadGirDirFile(name string) (string, error) {
	return util.ReadFileTrimNewline(g.GitDirPath(name))
}

func (g *GitRepo) BranchInfo(cfg config.BgpsConfig) (string, error) {
	var err error
	ref := ""
	step := ""
	total := ""

	if g.IsGitDir("rebase-merge") {
		if ref, err = g.ReadGirDirFile("rebase-merge/head-name"); err != nil {
			return "", err
		}
		if step, err = g.ReadGirDirFile("rebase-merge/msgnum"); err != nil {
			return "", err
		}
		if total, err = g.ReadGirDirFile("rebase-merge/end"); err != nil {
			return "", err
		}
		g.PromptMergeStatus = "|REBASE-m"
		if exists, exists_err := g.GitDirFileExists("rebase-merge/interactive"); exists {
			g.PromptMergeStatus = "|REBASE-i"
		} else if exists_err != nil {
			return "", exists_err
		}

	} else {
		if g.IsGitDir("rebase-apply") {

			step, err = g.ReadGirDirFile("rebase-apply/next")
			if err != nil {
				return "", err
			}
			total, err = g.ReadGirDirFile("rebase-apply/last")
			if err != nil {
				return "", err
			}
			rebasing, err := g.GitDirFileExists("rebase-apply/rebasing")
			if err != nil {
				return "", err
			}
			if rebasing {
				ref, err = g.ReadGirDirFile("rebase-apply/head-name")
				if err != nil {
					return "", err
				}
				g.PromptMergeStatus = "|REBASE"
			} else {
				// TODO: check if we need to get branch name here, bgps was not doing it
				applying, err := g.GitDirFileExists("rebase-apply/applying")
				if err != nil {
					return "", err
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
				if ref, err = SymbolicRef("HEAD"); err != nil {
					return "", err
				}
			} else {
				head := ""
				if head, err = g.ReadGirDirFile("HEAD"); err != nil {
					return "", err
				}
				ref = strings.TrimPrefix(head, "ref: ")

				if head == ref {
					tag, err := DescribeTag("HEAD")
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
		if _, err = LsFilesUnmerged(); err != nil {
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
	checkSparse := true
	if checkSparse {
		g.IsSparseCheckout, err = SparseCheckout()
		if err != nil {
			return "", err
		}

		if g.IsSparseCheckout {
			g.PromptSparseCheckoutStatus = "|SPARSE"
		}
	}

	prompt := fmt.Sprintf("%s%s%s%s", g.PromptBareRepoStatus, g.PromptBranch, g.PromptSparseCheckoutStatus, g.PromptMergeStatus)

	return prompt, nil
}

func (g *GitRepo) BranchStatus(cfg config.BgpsConfig) (string, string, error) {

	status := ""
	statusColor := ""

	if g.IsInBareRepo || g.IsInGitDir {
		c, _ := color.Color(strings.Split(cfg.ColorNoUpstream, " ")...)
		return status, c, nil
	}

	cleanWorkingTree, err := HasCleanWorkingTree()
	if err != nil {
		return "", "", err
	}
	hasUntracked, err := HasUntracked()
	if err != nil {
		return "", "", err
	}

	ahead, behind := 0, 0
	if g.Tag == "" && g.ShortSha != "" {
		ahead, behind, err = CommitCounts()
	}
	if err != nil {
		return "", "", err
	}

	if cleanWorkingTree {
		statusColor, _ = color.Color(strings.Split(cfg.ColorClean, " ")...)
	}

	if ahead > 0 {
		statusColor, _ = color.Color(strings.Split(cfg.ColorConflict, " ")...)
		status = fmt.Sprintf(cfg.AheadFormat, ahead)
	}
	if behind > 0 {
		statusColor, _ = color.Color(strings.Split(cfg.ColorConflict, " ")...)
		status = fmt.Sprintf(cfg.BehindFormat, behind)
	}

	if ahead > 0 && behind > 0 {
		status = fmt.Sprintf(cfg.DivergedFormat, ahead, behind)
	}

	if g.ShortSha == "" {
		statusColor, _ = color.Color(strings.Split(cfg.ColorNoUpstream, " ")...)
	}

	if hasUntracked {
		statusColor, _ = color.Color(strings.Split(cfg.ColorUntracked, " ")...)
		status = fmt.Sprintf("*%s", status)
	}

	if !cleanWorkingTree && !hasUntracked {
		statusColor, _ = color.Color(strings.Split(cfg.ColorDirty, " ")...)
		status = fmt.Sprintf("*%s", status)
	}
	if status != "" {
		status = " " + status
	}

	return status, statusColor, nil
}
