package globals

import "github.com/gdamore/tcell/v2"

var (
	// DefStyle ...
	DefStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	// CommentStyle
	CommentStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.Color18)
)
