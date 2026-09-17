package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// track.go draws everything around the buses: the exterior fences (vallas),
// the lane divider between the two buses, and the finish line (meta) at the
// end of the track. All of it is composed on top of a fixed-width canvas so
// it never overwrites the columns a bus is drawn on.

const (
	finishFlagWidth = 2 // columns occupied by the checkered finish line
	trackEdgeMargin = 1 // breathing room between the finish line and the terminal edge
	minTrackWidth   = 40
)

var (
	fencePostStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8D6E63"))
	fenceRailStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6D4C41"))
	laneLineStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFEB3B"))
	flagDarkStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#111111"))
	flagLightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA"))
)

const (
	flagDarkColor  = "#111111"
	flagLightColor = "#FAFAFA"
)

// flagColor returns the checkered finish-line color for this row, as a raw
// hex string so it can be used as a background fill behind a half-opaque
// sprite pixel (see busRowLine), not just rendered standalone.
func flagColor(row int) string {
	if row%2 == 1 {
		return flagLightColor
	}
	return flagDarkColor
}

// FinishColumn returns the column (0-indexed) where the finish line starts,
// given the total track width.
func FinishColumn(trackWidth int) int {
	col := trackWidth - finishFlagWidth - trackEdgeMargin
	if col < 0 {
		col = 0
	}
	return col
}

// TargetDistance returns how far (in columns) a bus must travel from column
// 0 for the frontmost *visible* pixel of its sprite (frontOffset) to land
// exactly on the finish line, stretching the race across the whole terminal.
func TargetDistance(trackWidth, frontOffset int) float64 {
	dist := FinishColumn(trackWidth) - frontOffset
	if dist < 1 {
		dist = 1
	}
	return float64(dist)
}

// flagGlyph renders a single checkered-flag column cell for this row.
func flagGlyph(row int) string {
	style := flagDarkStyle
	if row%2 == 1 {
		style = flagLightStyle
	}
	return style.Render("█")
}

func flagRun(row, width int) string {
	if width <= 0 {
		return ""
	}
	return strings.Repeat(flagGlyph(row), width)
}

// composeTrackLine pads a fully opaque `content` (whose visible width is
// contentWidth) up to the finish line and appends the finish-line cell for
// this row. Used by the fence and lane-divider lines, which never overlap
// the finish line themselves.
func composeTrackLine(content string, contentWidth, trackWidth, row int) string {
	finishCol := FinishColumn(trackWidth)
	if finishCol < contentWidth {
		finishCol = contentWidth
	}
	pad := finishCol - contentWidth
	return content + strings.Repeat(" ", pad) + flagRun(row, finishFlagWidth)
}

// fenceLine draws one exterior fence row (top or bottom valla) spanning up to
// the finish line, then lets composeTrackLine attach the finish-line cell.
func fenceLine(trackWidth, row int, top bool) string {
	span := FinishColumn(trackWidth)
	var b strings.Builder
	for i := 0; i < span; i++ {
		if i%4 == 0 {
			if top {
				b.WriteString(fencePostStyle.Render("▀"))
			} else {
				b.WriteString(fencePostStyle.Render("▄"))
			}
		} else {
			b.WriteString(fenceRailStyle.Render("─"))
		}
	}
	return composeTrackLine(b.String(), span, trackWidth, row)
}

// laneDividerLine draws the dashed lane-separator line between the two buses.
func laneDividerLine(trackWidth, row int) string {
	span := FinishColumn(trackWidth)
	var b strings.Builder
	for i := 0; i < span; i++ {
		if i%4 < 2 {
			b.WriteString(laneLineStyle.Render("-"))
		} else {
			b.WriteString(" ")
		}
	}
	return composeTrackLine(b.String(), span, trackWidth, row)
}

// busRowLine composites one terminal row of a bus sprite onto the track: the
// background is the checkered finish-line flag across its two columns, then
// spaces everywhere else. Each sprite cell is then composed against the
// background color at its column — a transparent half lets the finish flag
// (or space) show through it instead of being blanked out, including when
// only half of a half-block sprite cell is opaque.
func busRowLine(cells []SpriteCell, pos, row, trackWidth int) string {
	finishCol := FinishColumn(trackWidth)
	isFlag := func(idx int) bool {
		return idx >= finishCol && idx < finishCol+finishFlagWidth
	}

	base := make([]string, trackWidth)
	for i := range base {
		if isFlag(i) {
			base[i] = renderHalfBlock(flagColor(row), flagColor(row))
		} else {
			base[i] = " "
		}
	}

	for j, cell := range cells {
		idx := pos + j
		if idx < 0 || idx >= trackWidth {
			continue
		}
		if !cell.Opaque() {
			continue
		}
		top, bottom := cell.Top, cell.Bottom
		if isFlag(idx) {
			if top == "" {
				top = flagColor(row)
			}
			if bottom == "" {
				bottom = flagColor(row)
			}
		}
		base[idx] = renderHalfBlock(top, bottom)
	}

	return strings.Join(base, "")
}
