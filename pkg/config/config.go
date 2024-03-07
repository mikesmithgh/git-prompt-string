package config

type BgpsConfig struct {
	PromptPrefix   string `toml:"prompt_prefix"`
	PromptSuffix   string `toml:"prompt_suffix"`
	AheadFormat    string `toml:"ahead_format"`
	BehindFormat   string `toml:"behind_format"`
	DivergedFormat string `toml:"diverged_format"`
}
