package integration

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGitPromptString(t *testing.T) {
	tests := []struct {
		dir      string
		input    []string
		expected string
		environ  []string
		err      error
	}{
		{"bare", []string{"--config=NONE"}, "\x1b[90m \ue0a0 BARE:main\x1b[0m", nil, nil},
		{"no_upstream", []string{"--config=NONE"}, "\x1b[90m \ue0a0 main\x1b[0m", nil, nil},
		{"no_upstream_remote", []string{"--config=NONE"}, "\x1b[90m \ue0a0 main → mikesmithgh/test/main\x1b[0m", nil, nil},
		{"git_dir", []string{"--config=NONE"}, "\x1b[90m \ue0a0 GIT_DIR!\x1b[0m", nil, nil},
		{"clean", []string{"--config=NONE"}, "\x1b[32m \ue0a0 main\x1b[0m", nil, nil},
		{"tag", []string{"--config=NONE"}, "\x1b[90m \ue0a0 (v1.0.0)\x1b[0m", nil, nil},
		{"commit", []string{"--config=NONE"}, "\x1b[90m \ue0a0 (24afc95)\x1b[0m", nil, nil},
		{"dirty", []string{"--config=NONE"}, "\x1b[31m \ue0a0 main *\x1b[0m", nil, nil},
		{"dirty_staged", []string{"--config=NONE"}, "\x1b[31m \ue0a0 main *\x1b[0m", nil, nil},
		{"conflict_ahead", []string{"--config=NONE"}, "\x1b[33m \ue0a0 main ↑[1]\x1b[0m", nil, nil},
		{"conflict_behind", []string{"--config=NONE"}, "\x1b[33m \ue0a0 main ↓[1]\x1b[0m", nil, nil},
		{"conflict_diverged", []string{"--config=NONE"}, "\x1b[33m \ue0a0 main ↕ ↑[1] ↓[1]\x1b[0m", nil, nil},
		{"untracked", []string{"--config=NONE"}, "\x1b[35m \ue0a0 main *\x1b[0m", nil, nil},
		{"sparse", []string{"--config=NONE"}, "\x1b[32m \ue0a0 main|SPARSE\x1b[0m", nil, nil},
		{"sparse_merge_conflict", []string{"--config=NONE"}, "\x1b[31m \ue0a0 main|SPARSE|MERGING|CONFLICT *↕ ↑[1] ↓[1]\x1b[0m", nil, nil},

		// rebase merge
		{"rebase_i", []string{"--config=NONE"}, "\x1b[34m \ue0a0 main|REBASE-i 1/1\x1b[0m", nil, nil},
		{"rebase_m", []string{"--config=NONE"}, "\x1b[34m \ue0a0 main|REBASE-m 1/1\x1b[0m", nil, nil},
		// rebase apply
		{"am_rebase", []string{"--config=NONE"}, "\x1b[34m \ue0a0 (b69e688)|AM/REBASE 1/1\x1b[0m", nil, nil},
		{"am", []string{"--config=NONE"}, "\x1b[34m \ue0a0 (b69e688)|AM 1/1\x1b[0m", nil, nil},
		{"rebase", []string{"--config=NONE"}, "\x1b[34m \ue0a0 main|REBASE 1/1\x1b[0m", nil, nil},
		// merge
		{"merge_conflict", []string{"--config=NONE"}, "\x1b[31m \ue0a0 main|MERGING|CONFLICT *↕ ↑[1] ↓[1]\x1b[0m", nil, nil},
		{"merge", []string{"--config=NONE"}, "\x1b[35m \ue0a0 main|MERGING *↕ ↑[1] ↓[1]\x1b[0m", nil, nil},
		// cherry pick
		{"cherry_pick_conflict", []string{"--config=NONE"}, "\x1b[31m \ue0a0 main|CHERRY-PICKING|CONFLICT *↕ ↑[1] ↓[1]\x1b[0m", nil, nil},
		{"cherry_pick", []string{"--config=NONE"}, "\x1b[35m \ue0a0 main|CHERRY-PICKING *↕ ↑[1] ↓[1]\x1b[0m", nil, nil},
		// revert
		{"revert_conflict", []string{"--config=NONE"}, "\x1b[31m \ue0a0 main|REVERTING|CONFLICT *↕ ↑[2] ↓[1]\x1b[0m", nil, nil},
		{"revert", []string{"--config=NONE"}, "\x1b[31m \ue0a0 main|REVERTING *↕ ↑[2] ↓[1]\x1b[0m", nil, nil},
		// bisect
		{"bisect", []string{"--config=NONE"}, "\x1b[34m \ue0a0 main|BISECTING ↓[1]\x1b[0m", nil, nil},

		// formatting
		{"clean", []string{"--config=NONE", "--color-disabled"}, " \ue0a0 main", nil, nil},
		{"clean", []string{"--config=NONE", "--color-disabled", "--prompt-prefix= start "}, " start main", nil, nil},
		{"clean", []string{"--config=NONE", "--color-disabled", "--prompt-suffix= stop"}, " \ue0a0 main stop", nil, nil},
		{"conflict_ahead", []string{"--config=NONE", "--color-disabled", "--ahead-format=ahead by %d"}, " \ue0a0 main ahead by 1", nil, nil},
		{"conflict_behind", []string{"--config=NONE", "--color-disabled", "--behind-format=behind by %d"}, " \ue0a0 main behind by 1", nil, nil},
		{"conflict_diverged", []string{"--config=NONE", "--color-disabled", "--diverged-format=ahead by %d behind by %d"}, " \ue0a0 main ahead by 1 behind by 1", nil, nil},
		{"no_upstream_remote", []string{"--config=NONE", "--color-disabled", "--no-upstream-remote-format= upstream=[repo: %s branch: %s]"}, " \ue0a0 main upstream=[repo: mikesmithgh/test branch: main]", nil, nil},

		// color overrides
		{"clean", []string{"--config=../configs/color_overrides.toml"}, "\x1b[38;2;230;238;4m \ue0a0 main\x1b[0m", nil, nil},
		{"no_upstream", []string{"--config=../configs/color_overrides.toml"}, "\x1b[0m\x1b[30m\x1b[47m \ue0a0 main\x1b[0m", nil, nil},
		{"dirty", []string{"--config=../configs/color_overrides.toml"}, "\x1b[48;2;179;5;89m \ue0a0 main *\x1b[0m", nil, nil},
		{"conflict_ahead", []string{"--config=../configs/color_overrides.toml"}, "\x1b[38;2;252;183;40m \ue0a0 main ↑[1]\x1b[0m", nil, nil},
		{"untracked", []string{"--config=../configs/color_overrides.toml"}, "\x1b[38;2;255;0;0m\x1b[48;2;22;242;170m \ue0a0 main *\x1b[0m", nil, nil},
		{"bisect", []string{}, "\x1b[48;2;204;204;255m\x1b[35m \ue0a0 main|BISECTING ↓[1]\x1b[0m", []string{"GIT_PROMPT_STRING_CONFIG=../configs/color_overrides.toml"}, nil},

		// config errors
		{"clean", []string{"--config=/fromparam/does/not/exist"}, fmt.Sprintf("\x1b[31m git-prompt-string error(read config): \"open /fromparam/does/not/exist: %s\"\x1b[0m", notFoundMsg), nil, errors.New("exit status 1")},
		{"configs", []string{}, fmt.Sprintf("\x1b[31m git-prompt-string error(read config): \"open /fromenvvar/does/not/exist: %s\"\x1b[0m", notFoundMsg), []string{"GIT_PROMPT_STRING_CONFIG=/fromenvvar/does/not/exist"}, errors.New("exit status 1")},
		{"configs", []string{"--config=invalid_syntax.toml"}, fmt.Sprintf("\x1b[31m git-prompt-string error(unmarshal config): \"toml: expected character %s\"\x1b[0m", escapedEqualSign), nil, errors.New("exit status 1")},
		{"configs", []string{}, fmt.Sprintf("\x1b[31m git-prompt-string error(unmarshal config): \"toml: expected character %s\"\x1b[0m", escapedEqualSign), []string{"GIT_PROMPT_STRING_CONFIG=invalid_syntax.toml"}, errors.New("exit status 1")},

		{"norepo", []string{"--config=NONE"}, "", nil, nil},
	}

	for _, test := range tests {
		cmd := exec.Command(builtBinaryPath, test.input...)
		cmd.Dir = filepath.Join(tmpDir, "testdata", test.dir)
		if test.environ != nil {
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, test.environ...)
		}
		result, err := cmd.CombinedOutput()
		if test.err == nil && err != nil {
			t.Errorf("Unexpected error: %s", err)
		} else if test.err != nil && test.err.Error() != err.Error() {
			t.Errorf("Expected error: %s, got: %s", test.err, err)
		}
		actual := string(result)
		if actual != test.expected {
			t.Errorf("in directory %s, %s != %s\nexpected:\n%q, \ngot:\n%q", test.dir, test.expected, actual, test.expected, actual)
		}
	}
}
