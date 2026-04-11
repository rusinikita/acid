package terminal

import (
	"os"
	"runtime"
)

// SplitHint returns the keyboard shortcut to split a pane in the detected terminal.
func SplitHint() string {
	mod := "Cmd"
	if runtime.GOOS != "darwin" {
		mod = "Ctrl+Shift"
	}

	switch {
	case os.Getenv("TERM_PROGRAM") == "iTerm.app":
		return mod + "+D"
	case os.Getenv("GHOSTTY_RESOURCES_DIR") != "":
		return mod + "+D"
	case os.Getenv("KITTY_WINDOW_ID") != "":
		return "Ctrl+Shift+Enter"
	case os.Getenv("WEZTERM_PANE") != "":
		return mod + "+Shift+5"
	case os.Getenv("TMUX") != "":
		return `Ctrl+B then %`
	default:
		return "your terminal split shortcut"
	}
}
