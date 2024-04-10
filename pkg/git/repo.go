package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mikesmithgh/git-prompt-string/pkg/config"
	"github.com/mikesmithgh/git-prompt-string/pkg/util"
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

func (g *GitRepo) GitDirFileExistsExitOnError(name string) bool {
	exists, err := g.GitDirFileExists(name)
	if err != nil {
		util.ErrMsg(fmt.Sprintf("dir exists %s", name), err)
	}
	return exists
}

func (g *GitRepo) IsGitDir(name string) bool {
	return util.IsDir(g.GitDirPath(name))
}

func (g *GitRepo) IsGitDirSymlink(name string) bool {
	return util.IsSymlink(g.GitDirPath(name))
}

func (g *GitRepo) GitDirPath(path string) string {
	return filepath.Join(g.GitDir, path)
}

func (g *GitRepo) ReadGitDirFile(name string) (string, error) {
	return util.ReadFileTrimNewline(g.GitDirPath(name))
}

func (g *GitRepo) ReadGitDirFileExitOnError(name string) string {
	content, err := g.ReadGitDirFile(name)
	if err != nil {
		util.ErrMsg(fmt.Sprintf("read file %s", name), err)
	}
	return content
}

func (g *GitRepo) BranchInfo(cfg config.GitPromptStringConfig) (string, error) {
	var err error
	ref := ""
	step := ""
	total := ""

	if g.IsGitDir("rebase-merge") {
		ref = g.ReadGitDirFileExitOnError("rebase-merge/head-name")
		step = g.ReadGitDirFileExitOnError("rebase-merge/msgnum")
		total = g.ReadGitDirFileExitOnError("rebase-merge/end")
		g.PromptMergeStatus = "|REBASE-m"
		if g.GitDirFileExistsExitOnError("rebase-merge/interactive") {
			g.PromptMergeStatus = "|REBASE-i"
		}
	} else {
		switch {
		case g.IsGitDir("rebase-apply"):
			step = g.ReadGitDirFileExitOnError("rebase-apply/next")
			total = g.ReadGitDirFileExitOnError("rebase-apply/last")
			switch {
			case g.GitDirFileExistsExitOnError("rebase-apply/rebasing"):
				ref = g.ReadGitDirFileExitOnError("rebase-apply/head-name")
				g.PromptMergeStatus = "|REBASE"
			case g.GitDirFileExistsExitOnError("rebase-apply/applying"):
				g.PromptMergeStatus = "|AM"
			default:
				g.PromptMergeStatus = "|AM/REBASE"
			}
		case g.GitDirFileExistsExitOnError("MERGE_HEAD"):
			g.PromptMergeStatus = "|MERGING"
		case g.GitDirFileExistsExitOnError("CHERRY_PICK_HEAD"):
			g.PromptMergeStatus = "|CHERRY-PICKING"
		case g.GitDirFileExistsExitOnError("REVERT_HEAD"):
			g.PromptMergeStatus = "|REVERTING"
		case g.GitDirFileExistsExitOnError("BISECT_LOG"):
			g.PromptMergeStatus = "|BISECTING"
		}

		if ref == "" {
			if g.IsGitDirSymlink("HEAD") {
				if ref, err = SymbolicRef("HEAD"); err != nil {
					return "", err
				}
			} else {
				head := g.ReadGitDirFileExitOnError("HEAD")
				ref = strings.TrimPrefix(head, "ref: ")
				if head == ref {
					tag, err := DescribeTag("HEAD")
					switch {
					case err == nil:
						ref = tag
						g.Tag = ref
					case g.ShortSha == "" && len(head) > 7:
						ref = head[:7]
					default:
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
		unmerged, err := LsFilesUnmerged()
		if err != nil {
			return "", err
		}
		if unmerged != "" {
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

	g.IsSparseCheckout, err = SparseCheckout()
	if err != nil {
		return "", err
	}

	if g.IsSparseCheckout {
		g.PromptSparseCheckoutStatus = "|SPARSE"
	}

	if g.Tag == "" && g.ShortSha == "" && g.PromptMergeStatus == "" {
		branch_remote, err := BranchRemote(g.PromptBranch)
		var branch_merge string
		if err == nil {
			branch_merge, err = BranchMerge(g.PromptBranch)
		}
		if err == nil {
			remoteParts := strings.SplitN(branch_remote, ":", 2)
			if len(remoteParts) == 2 {
				branch_remote = strings.TrimSuffix(remoteParts[1], ".git")
			}

			if branch_merge != "" {
				g.PromptBranch += fmt.Sprintf(cfg.NoUpstreamRemoteFormat, branch_remote, strings.TrimPrefix(branch_merge, "refs/heads/"))
			}
		}
	}

	prompt := fmt.Sprintf("%s%s%s%s", g.PromptBareRepoStatus, g.PromptBranch, g.PromptSparseCheckoutStatus, g.PromptMergeStatus)

	return prompt, nil
}

func (g *GitRepo) BranchStatus(cfg config.GitPromptStringConfig) (string, string, error) {
	status := ""
	statusColor := ""

	if g.IsInBareRepo || g.IsInGitDir {
		return status, cfg.ColorNoUpstream, nil
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
		statusColor = cfg.ColorClean
	}

	if ahead > 0 {
		statusColor = cfg.ColorDelta
		status = fmt.Sprintf(cfg.AheadFormat, ahead)
	}
	if behind > 0 {
		statusColor = cfg.ColorDelta
		status = fmt.Sprintf(cfg.BehindFormat, behind)
	}

	if ahead > 0 && behind > 0 {
		status = fmt.Sprintf(cfg.DivergedFormat, ahead, behind)
	}

	if g.ShortSha == "" {
		statusColor = cfg.ColorNoUpstream
	}

	if g.PromptMergeStatus != "" {
		statusColor = cfg.ColorMerging
	}

	if hasUntracked {
		statusColor = cfg.ColorUntracked
		status = fmt.Sprintf("*%s", status)
	}

	if !cleanWorkingTree && !hasUntracked {
		statusColor = cfg.ColorDirty
		status = fmt.Sprintf("*%s", status)
	}
	if status != "" {
		status = " " + status
	}

	return status, statusColor, nil
}
