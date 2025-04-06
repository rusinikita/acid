package run

import (
	"fmt"
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

	db              string
	currentSequence int
}

func NewRunTable(runner runner, db string) tea.Model {
	data := &eventData{
		onlyStepsMode: true,
	}

	t := table.New().Data(data)

	m := &model{
		runner:  runner,
		running: true,

		table: t,
		data:  data,

		help: help.New(),
		db:   db,
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

		if key.Matches(msg, theme.DefaultKB.Mode) {
			m.data.onlyStepsMode = !m.data.onlyStepsMode
			m.vp.YOffset = 0
			m.UpdateViewport()

			return m, nil
		}

	case tea.WindowSizeMsg:
		m.window = msg

		m.UpdateViewport()
	case router.Message:
		m.data.clean()
		m.currentSequence = msg.DataInt

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
	msg := m.window

	m.table.Height(msg.Height)
	m.table.Width(msg.Width)
	m.vp.Width = msg.Width
	m.vp.Height = msg.Height - m.bordersHeight()

	m.vp.SetContent(m.table.String())
}

func (m *model) bordersHeight() int {
	return 1
}

func (m *model) View() string {
	menu := "  " + m.help.ShortHelpView(theme.DefaultKB.Menu())
	dbLabel := fmt.Sprint("seq ", m.currentSequence+1, " ", m.db)
	dbLabelLen := len(dbLabel)
	dbLabel = theme.StatusLineStyle.Render(dbLabel)
	fillWidth := m.window.Width - len(menu) - len(dbLabel)
	if fillWidth < 0 {
		fillWidth = 0
	}

	bottom := lipgloss.PlaceHorizontal(m.window.Width, lipgloss.Left, menu)
	bottom = bottom[:len(bottom)-dbLabelLen-4] + dbLabel

	return lipgloss.JoinVertical(lipgloss.Top, m.vp.View(), bottom)
}

type newEvent struct {
	Event       event.Event
	WaitingMore bool
}
