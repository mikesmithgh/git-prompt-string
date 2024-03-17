package integration

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBGPS(t *testing.T) {
	tests := []struct {
		dir      string
		input    []string
		expected string
		err      error
	}{
		{"bare", []string{"--config=NONE"}, "\x1b[90m \ue0a0 BARE:main\x1b[0m", nil},
		{"no_upstream", []string{"--config=NONE"}, "\x1b[90m \ue0a0 main\x1b[0m", nil},
		{"no_upstream_remote", []string{"--config=NONE"}, "\x1b[90m \ue0a0 main → mikesmithgh/test/main\x1b[0m", nil},
		{"git_dir", []string{"--config=NONE"}, "\x1b[90m \ue0a0 GIT_DIR!\x1b[0m", nil},
		{"clean", []string{"--config=NONE"}, "\x1b[32m \ue0a0 main\x1b[0m", nil},
		{"tag", []string{"--config=NONE"}, "\x1b[90m \ue0a0 (v1.0.0)\x1b[0m", nil},
		{"commit", []string{"--config=NONE"}, "\x1b[90m \ue0a0 (24afc95)\x1b[0m", nil},
		{"dirty", []string{"--config=NONE"}, "\x1b[31m \ue0a0 main *\x1b[0m", nil},
		{"dirty_staged", []string{"--config=NONE"}, "\x1b[31m \ue0a0 main *\x1b[0m", nil},
		{"conflict_ahead", []string{"--config=NONE"}, "\x1b[33m \ue0a0 main ↑[1]\x1b[0m", nil},
		{"conflict_behind", []string{"--config=NONE"}, "\x1b[33m \ue0a0 main ↓[1]\x1b[0m", nil},
		{"conflict_diverged", []string{"--config=NONE"}, "\x1b[33m \ue0a0 main ↕ ↑[1] ↓[1]\x1b[0m", nil},
		{"untracked", []string{"--config=NONE"}, "\x1b[35m \ue0a0 main *\x1b[0m", nil},
		{"sparse", []string{"--config=NONE"}, "\x1b[32m \ue0a0 main|SPARSE\x1b[0m", nil},
		{"sparse_merge_conflict", []string{"--config=NONE"}, "\x1b[31m \ue0a0 main|SPARSE|MERGING|CONFLICT *↕ ↑[1] ↓[1]\x1b[0m", nil},

		// rebase merge
		{"rebase_i", []string{"--config=NONE"}, "\x1b[34m \ue0a0 main|REBASE-i 1/1\x1b[0m", nil},
		{"rebase_m", []string{"--config=NONE"}, "\x1b[34m \ue0a0 main|REBASE-m 1/1\x1b[0m", nil},
		// rebase apply
		{"am_rebase", []string{"--config=NONE"}, "\x1b[34m \ue0a0 (b69e688)|AM/REBASE 1/1\x1b[0m", nil},
		{"am", []string{"--config=NONE"}, "\x1b[34m \ue0a0 (b69e688)|AM 1/1\x1b[0m", nil},
		{"rebase", []string{"--config=NONE"}, "\x1b[34m \ue0a0 main|REBASE 1/1\x1b[0m", nil},
		// merge
		{"merge_conflict", []string{"--config=NONE"}, "\x1b[31m \ue0a0 main|MERGING|CONFLICT *↕ ↑[1] ↓[1]\x1b[0m", nil},
		{"merge", []string{"--config=NONE"}, "\x1b[35m \ue0a0 main|MERGING *↕ ↑[1] ↓[1]\x1b[0m", nil},
		// cherry pick
		{"cherry_pick_conflict", []string{"--config=NONE"}, "\x1b[31m \ue0a0 main|CHERRY-PICKING|CONFLICT *↕ ↑[1] ↓[1]\x1b[0m", nil},
		{"cherry_pick", []string{"--config=NONE"}, "\x1b[35m \ue0a0 main|CHERRY-PICKING *↕ ↑[1] ↓[1]\x1b[0m", nil},
		// revert
		{"revert_conflict", []string{"--config=NONE"}, "\x1b[31m \ue0a0 main|REVERTING|CONFLICT *↕ ↑[2] ↓[1]\x1b[0m", nil},
		{"revert", []string{"--config=NONE"}, "\x1b[31m \ue0a0 main|REVERTING *↕ ↑[2] ↓[1]\x1b[0m", nil},
		// bisect
		{"bisect", []string{"--config=NONE"}, "\x1b[34m \ue0a0 main|BISECTING ↓[1]\x1b[0m", nil},

		{"clean", []string{"--config=NONE", "--color-enabled=false"}, " \ue0a0 main", nil},
		{"clean", []string{"--config=NONE", "--color-enabled=false", "--prompt-prefix= start "}, " start main", nil},
		{"clean", []string{"--config=NONE", "--color-enabled=false", "--prompt-suffix= stop"}, " \ue0a0 main stop", nil},
		{"conflict_ahead", []string{"--config=NONE", "--color-enabled=false", "--ahead-format=ahead by %d"}, " \ue0a0 main ahead by 1", nil},
		{"conflict_behind", []string{"--config=NONE", "--color-enabled=false", "--behind-format=behind by %d"}, " \ue0a0 main behind by 1", nil},
		{"conflict_diverged", []string{"--config=NONE", "--color-enabled=false", "--diverged-format=ahead by %d behind by %d"}, " \ue0a0 main ahead by 1 behind by 1", nil},
		{"no_upstream_remote", []string{"--config=NONE", "--color-enabled=false", "--no-upstream-remote-format= upstream=[repo: %s branch: %s]"}, " \ue0a0 main upstream=[repo: mikesmithgh/test branch: main]", nil},

		// TODO: add tests for color overrides

		// TODO: add test env var config doesn't exist
		// TODO: add test bad toml
		{"clean", []string{"--config=/does/not/exist"}, "\x1b[31m bgps error(read config): open /does/not/exist: no such file or directory\x1b[0m", nil},

		{"norepo", []string{"--config=NONE"}, "", nil},
	}

	for _, test := range tests {
		cmd := exec.Command(builtBinaryPath, test.input...)
		cmd.Dir = filepath.Join(tmpDir, "testdata", test.dir)
		result, err := cmd.CombinedOutput()
		if (err != nil && test.err == nil) || (err != nil && test.err != nil && err.Error() != test.err.Error()) {
			t.Errorf("Expected error: %v, got: %v", test.err, err)
		}
		actual := string(result)
		if actual != test.expected {
			t.Errorf("in directory %s, %s != %s\nexpected:\n%q, \ngot:\n%q", test.dir, test.expected, actual, test.expected, actual)
		}
	}
}
