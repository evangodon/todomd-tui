package ui

import lg "github.com/charmbracelet/lipgloss"

var (
	yellow = lg.Color("#f9e2af")
	blue   = lg.Color("#89b4fa")
	mauve  = lg.Color("#cba6f7")
	green  = lg.Color("#a6e3a1")
	red    = lg.Color("#f38ba8")
	white  = lg.Color("#ffffff")

	GreenText    = lg.NewStyle().Foreground(green).Render
	RedText      = lg.NewStyle().Foreground(red).Render
	BlueText     = lg.NewStyle().Foreground(blue).Render
	YellowText   = lg.NewStyle().Foreground(yellow).Render
	SelectedText = lg.NewStyle().Foreground(white).Bold(true).Render
	DimText      = lg.NewStyle().Faint(true).Render
)

var Color = struct {
	Yellow lg.Color
	Blue   lg.Color
	Mauve  lg.Color
	Green  lg.Color
	Red    lg.Color
	White  lg.Color
}{
	Yellow: lg.Color("#f9e2af"),
	Blue:   lg.Color("#89b4fa"),
	Mauve:  lg.Color("#cba6f7"),
	Green:  lg.Color("#a6e3a1"),
	Red:    lg.Color("#f38ba8"),
	White:  lg.Color("#ffffff"),
}
