package main

import "strings"

// busGrid is the pixel-art sprite for a bus, using two grid rows per
// rendered terminal line (top/bottom half-block trick).
var busGrid = []string{
	"...RRRRRRRRRR....",
	".RRRRRRRRRRRRRRR.",
	".CCWWWWWWWWWWWCC.",
	".SSSSSSSSSSSSSSS.",
	".BBBBBBBBBBBBBBB.",
	".BBKKKBBBBBKKKBB.",
	"...KGK.....KGK...",
	"...KKK.....KKK...",
}

var bluePalette = map[byte]string{
	'R': "#14213D", // roof
	'C': "#26C6DA", // windshield
	'W': "#1A1F2E", // side windows
	'S': "#64B5F6", // light stripe
	'B': "#1E88E5", // body
	'K': "#101010", // wheel (rim)
	'G': "#90A4AE", // wheel center
}

var redPalette = map[byte]string{
	'R': "#3D1414",
	'C': "#26C6DA",
	'W': "#1A1F2E",
	'S': "#F48FB1",
	'B': "#E53935",
	'K': "#101010",
	'G': "#90A4AE",
}

var pinkPalette = map[byte]string{
	'R': "#3D142E", // roof (dark pink/wine)
	'C': "#26C6DA", // windshield (cyan/glass)
	'W': "#1A1F2E", // side windows (dark blue)
	'S': "#F8BBD0", // light stripe (pastel pink)
	'B': "#E91E63", // body (intense pink)
	'K': "#101010", // wheel (rim)
	'G': "#90A4AE", // wheel center
}

// SpriteCell is one rendered terminal cell of a sprite, kept as the raw
// top/bottom half-block colors (rather than a pre-rendered glyph) so the
// final compositing can happen against whatever background sits behind each
// half at draw time — see busRowLine. An empty color means that half is
// transparent.
type SpriteCell struct {
	Top, Bottom string
}

// Opaque reports whether either half of the cell paints anything.
func (c SpriteCell) Opaque() bool {
	return c.Top != "" || c.Bottom != ""
}

// Sprite is a sprite pre-rendered into colored half-block cells, along with
// the geometry needed to place it and to know exactly which column its
// visible (opaque) artwork actually ends at.
type Sprite struct {
	Width, Height int
	// FrontOffset is the rightmost column (0-indexed, relative to the
	// sprite's own left edge) that contains an opaque pixel in any row.
	// Grids commonly have transparent padding columns on the right, so this
	// is what "the tip of the bus" really means, not Width-1.
	FrontOffset int
	Rows        [][]SpriteCell
}

// NewSprite renders `grid` with `palette` into a reusable Sprite. It is pure
// (depends only on its inputs), so it's safe to build once and reuse.
func NewSprite(grid []string, palette map[byte]string) Sprite {
	rows := normalizeGrid(grid)
	width := len(rows[0])

	sprite := Sprite{Width: width, Height: len(rows) / 2}
	sprite.Rows = make([][]SpriteCell, 0, sprite.Height)

	for i := 0; i+1 < len(rows); i += 2 {
		top, bottom := rows[i], rows[i+1]
		line := make([]SpriteCell, width)
		for j := 0; j < width; j++ {
			line[j] = renderCell(top[j], bottom[j], palette)
			if line[j].Opaque() && j > sprite.FrontOffset {
				sprite.FrontOffset = j
			}
		}
		sprite.Rows = append(sprite.Rows, line)
	}
	return sprite
}

func renderCell(t, b byte, palette map[byte]string) SpriteCell {
	cell := SpriteCell{}
	if !transparent(t) {
		cell.Top = palette[t]
	}
	if !transparent(b) {
		cell.Bottom = palette[b]
	}
	return cell
}

func normalizeGrid(grid []string) []string {
	maxLen := 0
	for _, row := range grid {
		if len(row) > maxLen {
			maxLen = len(row)
		}
	}

	out := make([]string, len(grid))
	for i, row := range grid {
		if len(row) < maxLen {
			row += strings.Repeat(".", maxLen-len(row))
		}
		out[i] = row
	}
	if len(out)%2 == 1 {
		out = append(out, strings.Repeat(".", maxLen))
	}
	return out
}

func transparent(c byte) bool {
	return c == '.' || c == ' '
}
