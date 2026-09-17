package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
)

const fps = 120

var tickInterval = time.Second / fps

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// raceModel is the bubbletea model driving the whole animation: physics,
// track sizing (based on the real terminal width) and rendering.
type raceModel struct {
	width, height int
	sized         bool

	blueBus, pinkbus Sprite

	spring harmonica.Spring
	target float64

	pos1, vel1 float64
	pos2, vel2 float64

	finished bool
	result   string
}

func newRaceModel() raceModel {
	return raceModel{
		spring: harmonica.NewSpring(harmonica.FPS(fps), 2.0, 2.0),
	}
}

func (m raceModel) Init() tea.Cmd {
	return nil
}

func (m raceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		trackWidth := m.width
		if trackWidth < minTrackWidth {
			trackWidth = minTrackWidth
		}

		firstResize := !m.sized
		if firstResize {
			m.sized = true
			m.blueBus = NewSprite(busGrid, bluePalette)
			m.pinkbus = NewSprite(busGrid, pinkPalette)
		}
		m.target = TargetDistance(trackWidth, m.blueBus.FrontOffset)

		if firstResize {
			return m, tick()
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tickMsg:
		if m.finished {
			return m, nil
		}
		m.step()
		if m.pos1 >= m.target || m.pos2 >= m.target {
			m.finished = true
			m.result = m.raceResult()
			return m, nil
		}
		return m, tick()
	}

	return m, nil
}

func (m *raceModel) step() {
	applyRandomImpulse(&m.vel1)
	applyRandomImpulse(&m.vel2)

	if m.pos1 < m.target {
		m.pos1, m.vel1 = m.spring.Update(m.pos1, m.vel1, m.target)
		if m.pos1 >= m.target {
			m.pos1, m.vel1 = m.target, 0
		}
	}
	if m.pos2 < m.target {
		m.pos2, m.vel2 = m.spring.Update(m.pos2, m.vel2, m.target)
		if m.pos2 >= m.target {
			m.pos2, m.vel2 = m.target, 0
		}
	}
}

func applyRandomImpulse(vel *float64) {
	if rand.Float64() >= 0.5 {
		return
	}
	impulse := rand.Float64() * 20.0
	if rand.Float64() < 0.35 {
		impulse = -impulse
	}
	*vel += impulse
	if *vel < 0 {
		*vel = 0
	}
}

func (m raceModel) raceResult() string {
	switch {
	case m.pos1 >= m.target && m.pos2 >= m.target:
		return "Tie!"
	case m.pos1 >= m.target:
		return "Blue bus wins!"
	default:
		return "Pink bus wins!"
	}
}

func (m raceModel) View() string {
	if !m.sized {
		return "Calculando el largo de la pista...\n"
	}

	trackWidth := m.width
	if trackWidth < minTrackWidth {
		trackWidth = minTrackWidth
	}

	var b strings.Builder
	row := 0

	writeLine := func(s string) {
		b.WriteString(s)
		b.WriteByte('\n')
		row++
	}

	writeLine(fenceLine(trackWidth, row, true))

	for _, cells := range m.blueBus.Rows {
		writeLine(busRowLine(cells, int(m.pos1), row, trackWidth))
	}

	writeLine(laneDividerLine(trackWidth, row))

	for _, cells := range m.pinkbus.Rows {
		writeLine(busRowLine(cells, int(m.pos2), row, trackWidth))
	}

	writeLine(fenceLine(trackWidth, row, false))

	if m.finished {
		b.WriteString("\n" + m.result + "  (q para salir)\n")
	} else {
		b.WriteString(fmt.Sprintf("\nPista: %d columnas · q para salir\n", trackWidth))
	}

	return b.String()
}
