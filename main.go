package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

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

func main() {
	fps := 120
	dt := time.Second / time.Duration(fps)

	spring := harmonica.NewSpring(harmonica.FPS(fps), 2.0, 2.0)

	pos1 := 0.0
	pos2 := 0.0
	vel1 := 0.0
	vel2 := 0.0
	targetPosition := 150.0

	fmt.Print("\033[2J")

	for {
		if rand.Float64() < 0.5 {
			impulse := rand.Float64() * 20.0
			if rand.Float64() < 0.35 {
				impulse = -impulse
			}
			vel1 += impulse
			if vel1 < 0 {
				vel1 = 0
			}
		}
		if rand.Float64() < 0.5 {
			impulse := rand.Float64() * 20.0
			if rand.Float64() < 0.35 {
				impulse = -impulse
			}
			vel2 += impulse
			if vel2 < 0 {
				vel2 = 0
			}
		}

		if pos1 < targetPosition {
			pos1, vel1 = spring.Update(pos1, vel1, targetPosition)
			if pos1 >= targetPosition {
				pos1 = targetPosition
				vel1 = 0
			}
		}

		if pos2 < targetPosition {
			pos2, vel2 = spring.Update(pos2, vel2, targetPosition)
			if pos2 >= targetPosition {
				pos2 = targetPosition
				vel2 = 0
			}
		}
		
		fmt.Print("\033[H")
		fmt.Print(RenderBus(busGrid, bluePalette, int(pos1)))
		fmt.Print("\n\n")
		fmt.Println(RenderBus(busGrid, redPalette, int(pos2)))

		if pos1 >= targetPosition || pos2 >= targetPosition {
			break
		}

		time.Sleep(dt)
	}

	switch {
	case pos1 >= targetPosition && pos2 >= targetPosition:
		fmt.Println("Tie!")
	case pos1 >= targetPosition:
		fmt.Println("Blue bus wins!")
	default:
		fmt.Println("Red bus wins!")
	}
}

func RenderBus(grid []string, palette map[byte]string, col int) string {
	rows := normalizeGrid(grid)
	indent := ""
	for i := 0; i < col; i++ {
		indent += " "
	}

	var builder strings.Builder
	for i := 0; i+1 < len(rows); i += 2 {
		// Add indentation at the beginning of each rendered row pair
		builder.WriteString(indent)
		builder.WriteString(renderPair(rows[i], rows[i+1], palette))
		builder.WriteString("\n")
	}

	return builder.String()
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

func renderPair(top, bottom string, palette map[byte]string) string {
	out := make([]byte, 0, len(top)*8)
	for i := 0; i < len(top); i++ {
		t, b := top[i], bottom[i]
		switch {
		case transparent(t) && transparent(b):
			out = append(out, ' ')
		case t == b:
			out = append(out, []byte(solid(t, palette))...)
		case transparent(t):
			out = append(out, []byte(lowerHalf(b, palette))...)
		case transparent(b):
			out = append(out, []byte(upperHalf(t, palette))...)
		default:
			out = append(out, []byte(split(t, b, palette))...)
		}
	}
	return string(out)
}

func solid(c byte, palette map[byte]string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(palette[c])).Render("█")
}

func upperHalf(c byte, palette map[byte]string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(palette[c])).Render("▀")
}

func lowerHalf(c byte, palette map[byte]string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(palette[c])).Render("▄")
}

func split(top, bottom byte, palette map[byte]string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(palette[top])).
		Background(lipgloss.Color(palette[bottom])).
		Render("▀")
}
