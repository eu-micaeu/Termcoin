package ui

import (
	"os"
)

const (
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
	ColorRed    = "\033[91m"
	ColorGreen  = "\033[92m"
	ColorYellow = "\033[93m"
	ColorBlue   = "\033[94m"
	ColorCyan   = "\033[96m"
	ColorGray   = "\033[90m"
)

var (
	Reset  = ColorReset
	Bold   = ColorBold
	Red    = ColorRed
	Green  = ColorGreen
	Yellow = ColorYellow
	Blue   = ColorBlue
	Cyan   = ColorCyan
	Gray   = ColorGray
)

func InitColors() {
	if !isTTY() {
		Reset = ""
		Bold = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Cyan = ""
		Gray = ""
	}
}

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
