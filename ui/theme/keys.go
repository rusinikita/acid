package theme

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var DefaultKB = KeyBindings{
	Up: key.NewBinding(
		key.WithKeys(tea.KeyUp.String(), "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys(tea.KeyDown.String(), "j"),
		key.WithHelp("↓/j", "down"),
	),
	Back: key.NewBinding(
		key.WithKeys(tea.KeyBackspace.String(), "b"),
		key.WithHelp("⌫/b", "back"),
	),
	Mode: key.NewBinding(
		key.WithKeys("m", tea.KeySpace.String()),
		key.WithHelp("m/space", "view mode"),
	),
	Learn: key.NewBinding(),
	ShowSetup: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "setup sql"),
	),
}

type KeyBindings struct {
	Up        key.Binding
	Down      key.Binding
	Back      key.Binding
	Mode      key.Binding
	Learn     key.Binding
	ShowSetup key.Binding
}

func (b KeyBindings) Menu() []key.Binding {
	return []key.Binding{
		b.Up,
		b.Down,
		b.Back,
		b.Mode,
		b.Learn,
		b.ShowSetup,
	}
}
