package config

type GPSConfig struct {
	PromptPrefix           string `toml:"prompt_prefix"`
	PromptSuffix           string `toml:"prompt_suffix"`
	AheadFormat            string `toml:"ahead_format"`
	BehindFormat           string `toml:"behind_format"`
	DivergedFormat         string `toml:"diverged_format"`
	NoUpstreamRemoteFormat string `toml:"no_upstream_remote_format"`
	ColorEnabled           bool   `toml:"color_enabled"`
	ColorClean             string `toml:"color_clean"`
	ColorConflict          string `toml:"color_conflict"`
	ColorDirty             string `toml:"color_dirty"`
	ColorUntracked         string `toml:"color_untracked"`
	ColorNoUpstream        string `toml:"color_no_upstream"`
	ColorMerging           string `toml:"color_merging"`
}
