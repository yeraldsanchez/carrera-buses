package main

import "github.com/charmbracelet/lipgloss"

// pixel.go composes a terminal cell out of a top and a bottom half-block
// color. Either half may be empty, meaning nothing paints there and the cell
// renders fully transparent (a plain space) — callers resolve transparency
// against whatever background belongs behind the cell *before* calling this,
// so a half-opaque sprite pixel can still show that background through its
// empty half instead of leaving it unset.
func renderHalfBlock(top, bottom string) string {
	switch {
	case top == "" && bottom == "":
		return " "
	case top == bottom:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(top)).Render("█")
	case bottom == "":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(top)).Render("▀")
	case top == "":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(bottom)).Render("▄")
	default:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(top)).
			Background(lipgloss.Color(bottom)).
			Render("▀")
	}
}
