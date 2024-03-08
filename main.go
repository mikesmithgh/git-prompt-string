package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/mikesmithgh/bgps/pkg/color"
	"github.com/mikesmithgh/bgps/pkg/config"
	"github.com/mikesmithgh/bgps/pkg/git"
	"github.com/mikesmithgh/bgps/pkg/util"
	"github.com/pelletier/go-toml/v2"
)

var (
	configPath      = flag.String("config", "", "")
	promptPrefix    = flag.String("prompt-prefix", "  ", "")
	promptSuffix    = flag.String("prompt-suffix", "", "")
	aheadFormat     = flag.String("ahead-format", "↑[%d]", "")
	behindFormat    = flag.String("behind-format", "↓[%d]", "")
	divergedFormat  = flag.String("diverged-format", "↕ ↑[%d] ↓[%d]", "")
	colorEnabled    = flag.Bool("color-enabled", true, "")
	colorClean      = flag.String("color-clean", "green", "")
	colorConflict   = flag.String("color-conflict", "yellow", "")
	colorDirty      = flag.String("color-dirty", "red", "")
	colorUntracked  = flag.String("color-untracked", "magenta", "")
	colorNoUpstream = flag.String("color-no-upstream", "bright-black", "")
)

func main() {

	cfg := config.BgpsConfig{
		PromptPrefix:    *promptPrefix,
		PromptSuffix:    *promptSuffix,
		AheadFormat:     *aheadFormat,
		BehindFormat:    *behindFormat,
		DivergedFormat:  *divergedFormat,
		ColorEnabled:    *colorEnabled,
		ColorClean:      *colorClean,
		ColorConflict:   *colorConflict,
		ColorDirty:      *colorDirty,
		ColorUntracked:  *colorUntracked,
		ColorNoUpstream: *colorNoUpstream,
	}

	flag.Parse()

	var bgpsConfigPath string
	if *configPath == "" {
		bgpsConfigPath = os.Getenv("BGPS_CONFIG")
	} else {
		bgpsConfigPath = *configPath
	}
	if bgpsConfigPath == "" {
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
		bgpsConfigPath = path.Join(xdgConfigHome, "bgps", "config.toml")
	}

	if bgpsConfigPath != "NONE" {
		bgpsConfigRaw, err := os.ReadFile(bgpsConfigPath)
		if err != nil && !os.IsNotExist(err) {
			util.ErrMsg("read config", err, 0)
		}

		err = toml.Unmarshal(bgpsConfigRaw, &cfg)
		if err != nil {
			util.ErrMsg("unmarshal config", err, 0)
		}
	}

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "prompt-prefix":
			cfg.PromptPrefix = f.Value.String()
		case "prompt-suffix":
			cfg.PromptSuffix = f.Value.String()
		case "ahead-format":
			cfg.AheadFormat = f.Value.String()
		case "behind-format":
			cfg.BehindFormat = f.Value.String()
		case "diverged-format":
			cfg.DivergedFormat = f.Value.String()
		case "color-enabled":
			cfg.ColorEnabled = f.Value.String() == f.DefValue
		case "color-clean":
			cfg.ColorClean = f.Value.String()
		case "color-conflict":
			cfg.ColorConflict = f.Value.String()
		case "color-dirty":
			cfg.ColorDirty = f.Value.String()
		case "color-untracked":
			cfg.ColorUntracked = f.Value.String()
		case "color-no-upstream":
			cfg.ColorUntracked = f.Value.String()
		}
	})

	if !cfg.ColorEnabled {
		color.Disable()
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
