package run

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/rusinikita/acid/event"
	"github.com/rusinikita/acid/sequence"
	"github.com/rusinikita/acid/ui/router"
	"github.com/rusinikita/acid/ui/theme"
)

type model struct {
	runner runner

	running bool
	window  tea.WindowSizeMsg

	table *table.Table
	vp    viewport.Model
	data  *eventData

	help help.Model
}

func NewRunTable(runner runner) tea.Model {
	data := &eventData{}

	t := table.New().Data(data)

	m := &model{
		runner:  runner,
		running: true,

		table: t,
		data:  data,

		help: help.New(),
	}

	styleFunc := func(row, _ int) lipgloss.Style {
		switch {
		case row == table.HeaderRow:
			return theme.HeaderStyle
		default:
			return theme.OddRowStyle.Width(m.window.Width/m.data.Columns() - m.data.Columns())
		}
	}
	t.StyleFunc(styleFunc)

	return m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func readNextEvent(r runner) func() tea.Msg {
	return func() tea.Msg {
		e, hasMore := r.Next()

		return newEvent{
			Event:       e,
			WaitingMore: hasMore,
		}
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.vp, cmd = m.vp.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, theme.DefaultKB.Back) {
			return m, router.Route("list")
		}
	case tea.WindowSizeMsg:
		m.table.Height(msg.Height)
		m.table.Width(msg.Width)
		m.vp.Width = msg.Width
		m.vp.Height = msg.Height - 1

		m.window = msg

		m.UpdateViewport()
	case router.Message:
		m.data.clean()

		run := func() tea.Msg {
			m.runner.Run(sequence.Sequences[msg.DataInt])

			return nil
		}

		return m, tea.Sequence(run, readNextEvent(m.runner))
	case newEvent:
		m.running = msg.WaitingMore
		if !m.running {
			return m, nil
		}

		m.data.add(msg.Event)
		m.table.Headers(m.data.headers()...)

		m.UpdateViewport()

		return m, readNextEvent(m.runner)
	}

	return m, cmd
}

func (m *model) UpdateViewport() {
	m.vp.SetContent(m.table.String())
}

func (m *model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, m.vp.View(), "  "+m.help.ShortHelpView(theme.DefaultKB.Menu()))
}

type newEvent struct {
	Event       event.Event
	WaitingMore bool
}
