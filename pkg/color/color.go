package color

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	esc = "\x1b"
)

var enabled bool = true

func codeToEscapeSequence(n int) string {
	return fmt.Sprintf("%s[%dm", esc, n)
}

var standardColors = map[string]string{
	"black":   codeToEscapeSequence(30),
	"red":     codeToEscapeSequence(31),
	"green":   codeToEscapeSequence(32),
	"yellow":  codeToEscapeSequence(33),
	"blue":    codeToEscapeSequence(34),
	"magenta": codeToEscapeSequence(35),
	"cyan":    codeToEscapeSequence(36),
	"white":   codeToEscapeSequence(37),

	"fg:black":   codeToEscapeSequence(30),
	"fg:red":     codeToEscapeSequence(31),
	"fg:green":   codeToEscapeSequence(32),
	"fg:yellow":  codeToEscapeSequence(33),
	"fg:blue":    codeToEscapeSequence(34),
	"fg:magenta": codeToEscapeSequence(35),
	"fg:cyan":    codeToEscapeSequence(36),
	"fg:white":   codeToEscapeSequence(37),

	"bright-black":   codeToEscapeSequence(90),
	"bright-red":     codeToEscapeSequence(91),
	"bright-green":   codeToEscapeSequence(92),
	"bright-yellow":  codeToEscapeSequence(93),
	"bright-blue":    codeToEscapeSequence(94),
	"bright-magenta": codeToEscapeSequence(95),
	"bright-cyan":    codeToEscapeSequence(96),
	"bright-white":   codeToEscapeSequence(97),

	"fg:bright-black":   codeToEscapeSequence(90),
	"fg:bright-red":     codeToEscapeSequence(91),
	"fg:bright-green":   codeToEscapeSequence(92),
	"fg:bright-yellow":  codeToEscapeSequence(93),
	"fg:bright-blue":    codeToEscapeSequence(94),
	"fg:bright-magenta": codeToEscapeSequence(95),
	"fg:bright-cyan":    codeToEscapeSequence(96),
	"fg:bright-white":   codeToEscapeSequence(97),

	"bg:black":   codeToEscapeSequence(40),
	"bg:red":     codeToEscapeSequence(41),
	"bg:green":   codeToEscapeSequence(42),
	"bg:yellow":  codeToEscapeSequence(43),
	"bg:blue":    codeToEscapeSequence(44),
	"bg:magenta": codeToEscapeSequence(45),
	"bg:cyan":    codeToEscapeSequence(46),
	"bg:white":   codeToEscapeSequence(47),

	"bg:bright-black":   codeToEscapeSequence(100),
	"bg:bright-red":     codeToEscapeSequence(101),
	"bg:bright-green":   codeToEscapeSequence(102),
	"bg:bright-yellow":  codeToEscapeSequence(103),
	"bg:bright-blue":    codeToEscapeSequence(104),
	"bg:bright-magenta": codeToEscapeSequence(105),
	"bg:bright-cyan":    codeToEscapeSequence(106),
	"bg:bright-white":   codeToEscapeSequence(107),

	"none": codeToEscapeSequence(0),
}

func hexToRGB(hex string) (int, int, int, error) {
	if !strings.HasPrefix(hex, "#") {
		return 0, 0, 0, fmt.Errorf("hex must start with #, got %s", hex)
	}

	hex = strings.TrimPrefix(hex, "#")

	if len(hex) != 6 {
		return 0, 0, 0, fmt.Errorf("hex must be 6 digits, got %s", hex)
	}
	// Parse the hex string into RGB components
	r, err := strconv.ParseInt(hex[0:2], 16, 32)
	if err != nil {
		return 0, 0, 0, err
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 32)
	if err != nil {
		return 0, 0, 0, err
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 32)
	if err != nil {
		return 0, 0, 0, err
	}

	return int(r), int(g), int(b), nil
}

func rgbToEscapeSequence(r, g, b int, isBg bool) string {
	var colorType string
	if isBg {
		colorType = "48"
	} else {
		colorType = "38"
	}
	return fmt.Sprintf("\x1b[%s;2;%d;%d;%dm", colorType, r, g, b)
}

func Disable() {
	enabled = false
}

func Color(colors ...string) (string, error) {
	seq := ""
	if !enabled {
		return seq, nil
	}
	for _, color := range colors {
		if strings.HasPrefix(color, "#") || strings.HasPrefix(color, "fg:#") || strings.HasPrefix(color, "bg:#") {
			r, g, b, err := hexToRGB(strings.TrimPrefix(strings.TrimPrefix(color, "fg:"), "bg:"))
			if err != nil {
				return "", err
			}
			seq += rgbToEscapeSequence(r, g, b, strings.HasPrefix(color, "bg:#"))
		} else {
			s, exists := standardColors[color]
			if !exists {
				return "", fmt.Errorf("color %s not found", color)
			}
			seq += s
		}
	}
	return seq, nil
}
