package main

import (
	"github.com/charmbracelet/lipgloss"
)

// menu.go draws the floating, bordered menu/status box that sits below the
// track. The box itself has a clean, solid interior (border + padding around
// centered content) — the decorative "noise" texture of CJK glyphs lives
// entirely *outside* the box, filling the surrounding whitespace, the way
// Charm's marketing site cards float over a textured page background.

var (
	menuBorderColor   = lipgloss.Color("#E91E63")
	patternGlyphChars = "金継ぎ"
	patternColor      = lipgloss.Color("#1e1b3a")
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#fac0b6"))
	hintStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#5C5A72"))
	winnerPinkStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#E91E63"))
	winnerBlueStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#1E88E5"))
	statusStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#C9C6DC"))
	buttonOnStyle     = lipgloss.NewStyle().Bold(true).
				Foreground(lipgloss.Color("#0B0B14")).
				Background(lipgloss.Color("#B388FF")).
				Padding(0, 2)
	buttonOffStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8A8A9E")).Padding(0, 2)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(menuBorderColor).
			Padding(1, 3)
)

func renderButton(label string, selected bool) string {
	if selected {
		return buttonOnStyle.Render(label)
	}
	return buttonOffStyle.Render(label)
}

// boxContent builds the (unbordered) centered content for the model's
// current state: the start menu, an in-progress status, or the finish
// announcement.
func boxContent(m raceModel) string {
	switch m.state {
	case stateRacing:
		return statusStyle.Render("¡La carrera está en marcha!")
	case stateFinished:
		if m.result == "¡El bus azul ha ganado! 🥵 🥵 🥵" {
			return winnerBlueStyle.Render(m.result)
		}
		return winnerPinkStyle.Render(m.result)
	default: // stateMenu
		title := titleStyle.Render("CARRERA DE BUSES 🗣 🗣 🗣")
		buttons := lipgloss.JoinHorizontal(lipgloss.Center,
			renderButton("Iniciar", m.selected == 0),
			"  ",
			renderButton("Salir", m.selected == 1),
		)
		return lipgloss.JoinVertical(lipgloss.Center, title, "", buttons)
	}
}

// renderMenuBox lays the bordered card out on a trackWidth-wide canvas whose
// margins (one row above, one row below, and the sides) are filled with the
// decorative glyph pattern — everything inside the border stays a clean,
// solid panel.
func renderMenuBox(m raceModel, trackWidth int) string {
	box := boxStyle.Render(boxContent(m))
	placeHeight := lipgloss.Height(box) + 2

	return lipgloss.Place(trackWidth, placeHeight, lipgloss.Center, lipgloss.Center, box,
		lipgloss.WithWhitespaceChars(patternGlyphChars),
		lipgloss.WithWhitespaceForeground(patternColor),
	)
}

func renderHint(m raceModel) string {
	switch m.state {
	case stateMenu:
		return hintStyle.Render("←/→ elegir · enter confirmar · q salir")
	default:
		return hintStyle.Render("q para salir")
	}
}
