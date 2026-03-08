package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up          key.Binding
	Down        key.Binding
	Kill        key.Binding
	Filter      key.Binding
	ClearFilter key.Binding
	Copy        key.Binding
	Refresh     key.Binding
	Help        key.Binding
	Quit        key.Binding
	ForceKill   key.Binding
	Sort        key.Binding
	Open        key.Binding
	ToggleUDP   key.Binding
	ToggleDedup key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Kill: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "kill process"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	ClearFilter: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear filter"),
	),
	Copy: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy to clipboard"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh now"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	ForceKill: key.NewBinding(
		key.WithKeys("X"),
		key.WithHelp("X", "force kill (SIGKILL)"),
	),
	Sort: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "cycle sort"),
	),
	Open: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open in browser"),
	),
	ToggleUDP: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "toggle TCP/UDP"),
	),
	ToggleDedup: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "toggle merge"),
	),
}
