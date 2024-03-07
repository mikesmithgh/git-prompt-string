package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/mikesmithgh/bgps/pkg/config"
	"github.com/mikesmithgh/bgps/pkg/git"
	"github.com/mikesmithgh/bgps/pkg/util"
	"github.com/pelletier/go-toml/v2"
)

func main() {

	cfg := config.BgpsConfig{
		PromptPrefix:   "  ",
		PromptSuffix:   "",
		AheadFormat:    "↑[%d]",
		BehindFormat:   "↓[%d]",
		DivergedFormat: "↕ ↑[%d] ↓[%d]",
	}

	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			util.ErrMsg("user home", err, 0)
		}
		if runtime.GOOS == "windows" {
			xdgConfigHome = path.Join(home, "AppData", "Local")
		} else {
			xdgConfigHome = path.Join(home, ".config")
		}
	}
	bgpsConfigPath := path.Join(xdgConfigHome, "bgps.toml")

	bgpsConfigRaw, err := os.ReadFile(bgpsConfigPath)
	if err != nil && !os.IsNotExist(err) {
		util.ErrMsg("read config", err, 0)
	}

	err = toml.Unmarshal(bgpsConfigRaw, &cfg)
	if err != nil {
		util.ErrMsg("unmarshal config", err, 0)
	}

	gitRepo, stderr, err := git.RevParse()
	if err != nil {
		if strings.Contains(string(stderr), "not a git repository") {
			os.Exit(0)
		}
		// allow other errors to pass through, the git repo may not have upstream
	}

	branchInfo, err := gitRepo.BranchInfo(cfg)
	if err != nil {
		util.ErrMsg("branch info", err, 0)
	}
	branchStatus, color, err := gitRepo.BranchStatus(cfg)
	if err != nil {
		util.ErrMsg("branch status", err, 0)
	}

	fmt.Printf("%s%s%s%s%s", color, cfg.PromptPrefix, branchInfo, branchStatus, cfg.PromptSuffix)
}
