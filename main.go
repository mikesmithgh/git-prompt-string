package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
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
	configPath             = flag.String("config", "", "The filepath of the git-prompt-string toml configuration.")
	promptPrefix           = flag.String("prompt-prefix", " \ue0a0 ", "A prefix that is added to the beginning of the prompt. The\npowerline icon  is used be default. It is recommended to\nuse a Nerd Font to properly display the  (nf-pl-branch) icon.\nSee https://www.nerdfonts.com/ to download a Nerd Font. If you\ndo not want this symbol, replace the prompt prefix with \" \".\n\\ue0a0 is the unicode representation of .")
	promptSuffix           = flag.String("prompt-suffix", "", "A suffix that is added to the end of the prompt.")
	aheadFormat            = flag.String("ahead-format", "↑[%v]", "The format used to indicate the number of commits ahead of the\nremote branch. The %v verb represents the number of commits\nahead. One %v verb is required.")
	behindFormat           = flag.String("behind-format", "↓[%v]", "The format used to indicate the number of commits behind the\nremote branch. The %v verb represents the number of commits\nbehind. One %v verb is required.")
	divergedFormat         = flag.String("diverged-format", "↕ ↑[%v] ↓[%v]", "The format used to indicate the number of commits diverged\nfrom the remote branch. The first %v verb represents the number\nof commits ahead of the remote branch. The second %v verb\nrepresents the number of commits behind the remote branch. Two\n%v verbs are required.")
	noUpstreamRemoteFormat = flag.String("no-upstream-remote-format", " → %v/%v", "The format used to indicate when there is no remote upstream,\nbut there is still a remote branch configured. The first %v\nrepresents the remote repository. The second %v represents the\nremote branch. Two %v are required.")
	colorDisabled          = flag.Bool("color-disabled", false, "Disable all colors in the prompt.")
	colorClean             = flag.String("color-clean", "green", "The color of the prompt when the working directory is clean.\n")
	colorDelta             = flag.String("color-delta", "yellow", "The color of the prompt when the local branch is ahead, behind,\nor has diverged from the remote branch.")
	colorDirty             = flag.String("color-dirty", "red", "The color of the prompt when the working directory has changes\nthat have not yet been commited.")
	colorUntracked         = flag.String("color-untracked", "magenta", "The color of the prompt when there are untracked files in the\nworking directory.")
	colorNoUpstream        = flag.String("color-no-upstream", "bright-black", "The color of the prompt when there is no remote upstream branch.\n")
	colorMerging           = flag.String("color-merging", "blue", "The color of the prompt during a merge, rebase, cherry-pick,\nrevert, or bisect.")
	versionFlag            = flag.Bool("version", false, "Print version information for git-prompt-string.")
)

func header() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString("git-prompt-string: a shell agnostic git prompt written in Go.\n")
	sb.WriteString("https://github.com/mikesmithgh/git-prompt-string\n")
	sb.WriteString("\n")
	return sb.String()
}

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
		ColorDelta:             *colorDelta,
		ColorDirty:             *colorDirty,
		ColorUntracked:         *colorUntracked,
		ColorNoUpstream:        *colorNoUpstream,
		ColorMerging:           *colorMerging,
	}

	flag.Usage = func() {
		w := flag.CommandLine.Output()

		var sb strings.Builder

		sb.WriteString(header())
		sb.WriteString("Usage:")
		sb.WriteString("\n")
		sb.WriteString("git-prompt-string [flags]")
		sb.WriteString("\n\n")
		sb.WriteString("Flags can be prefixed with either - or --. For example, -version and")
		sb.WriteString("\n")
		sb.WriteString("--version are both valid flags.")
		sb.WriteString("\n\n")
		sb.WriteString("Flags:")
		sb.WriteString("\n")
		fmt.Fprint(w, sb.String())

		flag.PrintDefaults()
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
				util.ErrMsg("user home", err)
			}
			xdgConfigHome = path.Join(home, util.XDGConfigPath)
		}
		gpsConfig = path.Join(xdgConfigHome, "git-prompt-string", "config.toml")
	}

	if gpsConfig != "NONE" {
		gpsConfigRaw, err := os.ReadFile(gpsConfig)
		if err != nil && !os.IsNotExist(err) {
			util.ErrMsg("read config exists", err)
		}

		if err != nil && (*configPath != "" || gpsConfigEnv != "") {
			util.ErrMsg("read config", err)
		}

		err = toml.Unmarshal(gpsConfigRaw, &cfg)
		if err != nil {
			util.ErrMsg("unmarshal config", err)
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
				util.ErrMsg("parse color disabled", err)
			}
			cfg.ColorDisabled = colorDisabled
		case "color-clean":
			cfg.ColorClean = f.Value.String()
		case "color-delta":
			cfg.ColorDelta = f.Value.String()
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
		fmt.Print(header())
		fmt.Printf("Version:   %s\n", version)
		fmt.Printf("Commit:    %s\n", commit)
		fmt.Printf("BuildDate: %s\n", date)
		os.Exit(0)
	}

	clearColor, err := color.Color("none")
	if err != nil {
		util.ErrMsg("color none", err)
	}

	gitRepo, stderr, err := git.RevParse()
	if err != nil {
		switch {
		case strings.Contains(err.Error(), exec.ErrNotFound.Error()):
			util.ErrMsg("rev parse", err)
		case strings.Contains(string(stderr), "not a git repository"):
			os.Exit(0)
		default:
			// allow other errors to pass through, the git repo may not have upstream
		}
	}

	branchInfo, err := gitRepo.BranchInfo(cfg)
	if err != nil {
		util.ErrMsg("branch info", err)
	}
	branchStatus, promptColor, err := gitRepo.BranchStatus(cfg)
	if err != nil {
		util.ErrMsg("branch status", err)
	}

	fmt.Printf("%s%s%s%s%s%s", promptColor, cfg.PromptPrefix, branchInfo, branchStatus, cfg.PromptSuffix, clearColor)
}
