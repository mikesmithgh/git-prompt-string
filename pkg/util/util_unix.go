//go:build !windows

package util

import "path"

var (
	XDGConfigPath string = path.Join(".config")
)
