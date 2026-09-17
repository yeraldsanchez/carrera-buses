package main

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

const fps = 120

var tickInterval = time.Second / fps

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// raceResetMsg fires a few seconds after a race finishes, returning the menu
// box from the winner announcement back to the normal start/exit menu.
type raceResetMsg time.Time

const resultDisplayTime = 3 * time.Second

func resetAfterResult() tea.Cmd {
	return tea.Tick(resultDisplayTime, func(t time.Time) tea.Msg {
		return raceResetMsg(t)
	})
}

// appState is which screen the menu box (and the input handling) is in.
type appState int

const (
	stateMenu appState = iota
	stateRacing
	stateFinished
)

// raceModel is the bubbletea model driving the whole animation: physics,
// track sizing (based on the real terminal width) and rendering.
type raceModel struct {
	width, height int
	sized         bool

	state    appState
	selected int // 0 = Iniciar, 1 = Salir, only meaningful in stateMenu

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
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		switch m.state {
		case stateMenu:
			switch msg.String() {
			case "q", "esc":
				return m, tea.Quit
			case "left", "right", "h", "l", "tab":
				m.selected = 1 - m.selected
			case "enter", " ":
				if m.selected == 1 {
					return m, tea.Quit
				}
				m.pos1, m.vel1 = 0, 0
				m.pos2, m.vel2 = 0, 0
				m.finished = false
				m.result = ""
				m.state = stateRacing
				return m, tick()
			}
		case stateRacing, stateFinished:
			switch msg.String() {
			case "q", "esc":
				return m, tea.Quit
			}
		}

	case tickMsg:
		if m.state != stateRacing {
			return m, nil
		}
		m.step()
		if m.pos1 >= m.target || m.pos2 >= m.target {
			m.finished = true
			m.result = m.raceResult()
			m.state = stateFinished
			return m, resetAfterResult()
		}
		return m, tick()

	case raceResetMsg:
		if m.state == stateFinished {
			m.state = stateMenu
			m.selected = 0
			m.finished = false
			m.pos1, m.vel1 = 0, 0
			m.pos2, m.vel2 = 0, 0
		}
		return m, nil
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
		return "¡Empate!"
	case m.pos1 >= m.target:
		return "¡El bus azul ha ganado! 🥵 🥵 🥵"
	default:
		return "¡El bus rosa ha ganado! 🤪 🤪 🤪"
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

	b.WriteString("\n")
	b.WriteString(renderMenuBox(m, trackWidth))
	b.WriteString("\n")
	b.WriteString(lipgloss.PlaceHorizontal(trackWidth, lipgloss.Center, renderHint(m)))
	b.WriteString("\n")

	return b.String()
}
