# CarreraBuses

A bus racing game in the terminal (TUI) built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
and [Lip Gloss](https://github.com/charmbracelet/lipgloss). Two pixel-art buses
compete using spring physics on the track to the finish line; a menu
lets you start the race or exit, and when finished displays the winner.

## Requirements

- Go 1.27 or higher
- A terminal with color support and Unicode/wide characters (CJK)

### Recommended Terminals

- **macOS**: iTerm2, Kitty, Alacritty
- **Linux**: GNOME Terminal, Konsole, Kitty, Alacritty
- **Windows**: Windows Terminal, Kitty, Alacritty
- **Cross-platform**: Hyper, Alacritty

## How to Run

```bash
go run .
```

Or build the binary and run it:

```bash
go build -o CarreraBuses .
./CarreraBuses
```

Or download a prebuilt binary for your platform from the
[Releases page](https://github.com/yeraldsanchez/carrera-buses/releases).

## Controls

- `←` / `→`: choose between "Start" and "Exit"
- `Enter`: confirm the selected option
- `q` / `Esc` / `Ctrl+C`: exit at any time
