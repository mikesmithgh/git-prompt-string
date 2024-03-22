package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/mikesmithgh/git-prompt-string/pkg/color"
	"github.com/mikesmithgh/git-prompt-string/pkg/config"
	"github.com/mikesmithgh/git-prompt-string/pkg/git"
	"github.com/mikesmithgh/git-prompt-string/pkg/util"
	"github.com/pelletier/go-toml/v2"
)

var (
	version                = "dev"     // populated by goreleaser
	commit                 = "none"    // populated by goreleaser
	date                   = "unknown" // populated by goreleaser
	configPath             = flag.String("config", "", "")
	promptPrefix           = flag.String("prompt-prefix", " \ue0a0 ", "")
	promptSuffix           = flag.String("prompt-suffix", "", "")
	aheadFormat            = flag.String("ahead-format", "↑[%d]", "")
	behindFormat           = flag.String("behind-format", "↓[%d]", "")
	divergedFormat         = flag.String("diverged-format", "↕ ↑[%d] ↓[%d]", "")
	noUpstreamRemoteFormat = flag.String("no-upstream-remote-format", " → %s/%s", "")
	colorDisabled          = flag.Bool("color-disabled", false, "disable all color in prompt string")
	colorClean             = flag.String("color-clean", "green", "")
	colorConflict          = flag.String("color-conflict", "yellow", "")
	colorDirty             = flag.String("color-dirty", "red", "")
	colorUntracked         = flag.String("color-untracked", "magenta", "")
	colorNoUpstream        = flag.String("color-no-upstream", "bright-black", "")
	colorMerging           = flag.String("color-merging", "blue", "")
	versionFlag            = flag.Bool("version", false, "version for git-prompt-string")
)

func main() {
	cfg := config.GPSConfig{
		PromptPrefix:           *promptPrefix,
		PromptSuffix:           *promptSuffix,
		AheadFormat:            *aheadFormat,
		BehindFormat:           *behindFormat,
		DivergedFormat:         *divergedFormat,
		NoUpstreamRemoteFormat: *noUpstreamRemoteFormat,
		ColorDisabled:          *colorDisabled,
		ColorClean:             *colorClean,
		ColorConflict:          *colorConflict,
		ColorDirty:             *colorDirty,
		ColorUntracked:         *colorUntracked,
		ColorNoUpstream:        *colorNoUpstream,
		ColorMerging:           *colorMerging,
	}

	flag.Parse()

	var gpsConfig string
	gpsConfigEnv := os.Getenv("GIT_PROMPT_STRING_CONFIG")
	if *configPath == "" {
		gpsConfig = gpsConfigEnv
	} else {
		gpsConfig = *configPath
	}
	if gpsConfig == "" {
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
		gpsConfig = path.Join(xdgConfigHome, "git-prompt-string", "config.toml")
	}

	if gpsConfig != "NONE" {
		gpsConfigRaw, err := os.ReadFile(gpsConfig)
		if err != nil && !os.IsNotExist(err) {
			util.ErrMsg("read config exists", err, 0)
		}

		if err != nil && (*configPath != "" || gpsConfigEnv != "") {
			util.ErrMsg("read config", err, 0)
		}

		err = toml.Unmarshal(gpsConfigRaw, &cfg)
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
		case "no-upstream-remote-format":
			cfg.NoUpstreamRemoteFormat = f.Value.String()
		case "color-disabled":
			colorDisabled, err := strconv.ParseBool(f.Value.String())
			if err != nil {
				util.ErrMsg("parse color disabled", err, 0)
			}
			cfg.ColorDisabled = colorDisabled
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
		case "color-merging":
			cfg.ColorMerging = f.Value.String()
		}
	})

	if cfg.ColorDisabled {
		color.Disable()
	}

	if *versionFlag {
		fmt.Println()
		fmt.Println("git-prompt-string")
		fmt.Println("https://github.com/mikesmithgh/git-prompt-string")
		fmt.Println()
		fmt.Printf("Version:   %s\n", version)
		fmt.Printf("Commit:    %s\n", commit)
		fmt.Printf("BuildDate: %s\n", date)
		os.Exit(0)
	}

	clearColor, err := color.Color("none")
	if err != nil {
		util.ErrMsg("color none", err, 0)
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
	branchStatus, promptColor, err := gitRepo.BranchStatus(cfg)
	if err != nil {
		util.ErrMsg("branch status", err, 0)
	}

	fmt.Printf("%s%s%s%s%s%s", promptColor, cfg.PromptPrefix, branchInfo, branchStatus, cfg.PromptSuffix, clearColor)
}
