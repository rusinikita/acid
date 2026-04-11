package run

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/event"
	"github.com/rusinikita/acid/sequence"
	"github.com/rusinikita/acid/ui/router"
	"github.com/rusinikita/acid/ui/theme"
)

type model struct {
	runner   runner
	toggleCh <-chan struct{}

	running    bool
	serverMode bool
	window     tea.WindowSizeMsg

	table *table.Table
	vp    viewport.Model
	data  *eventData

	help help.Model

	db              string
	currentSequence int
	sequences       []sequence.Sequence
}

func NewRunTable(runner runner, db string, sequences []sequence.Sequence) tea.Model {
	data := &eventData{
		onlyStepsMode: true,
	}

	t := table.New().Data(data)

	m := &model{
		runner:    runner,
		running:   true,
		sequences: sequences,

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

// NewServerRunTable creates a run model for server mode. Events are received
// from the network via source; no local DB connection or sequence selection is needed.
// toggleCh fires when the /toggle-mode endpoint is called, switching result visibility.
func NewServerRunTable(source runner, port string, toggleCh <-chan struct{}) tea.Model {
	data := &eventData{
		onlyStepsMode: true,
		transactions:  []call.TrxID{""},
	}

	t := table.New().Data(data)

	m := &model{
		runner:     source,
		toggleCh:   toggleCh,
		running:    true,
		serverMode: true,

		table: t,
		data:  data,

		help: help.New(),
		db:   ":" + port,
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
	if m.serverMode {
		return tea.Batch(readNextEvent(m.runner), waitForToggle(m.toggleCh))
	}
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

type toggleModeMsg struct{}

func waitForToggle(ch <-chan struct{}) tea.Cmd {
	return func() tea.Msg {
		<-ch
		return toggleModeMsg{}
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

		if !m.serverMode && key.Matches(msg, theme.DefaultKB.Mode) {
			m.data.onlyStepsMode = !m.data.onlyStepsMode
			m.vp.YOffset = 0
			m.UpdateViewport()

			return m, nil
		}

		if key.Matches(msg, theme.DefaultKB.ShowSetup) {
			m.data.showSetupEvents = !m.data.showSetupEvents
			m.vp.YOffset = 0
			m.UpdateViewport()

			return m, nil
		}

	case tea.WindowSizeMsg:
		m.window = msg

		m.UpdateViewport()
	case router.Message:
		if m.serverMode {
			return m, nil
		}
		m.data.clean()
		m.currentSequence = msg.DataInt

		run := func() tea.Msg {
			m.runner.Run(m.sequences[msg.DataInt])

			return nil
		}

		return m, tea.Sequence(run, readNextEvent(m.runner))
	case toggleModeMsg:
		m.data.onlyStepsMode = !m.data.onlyStepsMode
		m.vp.YOffset = 0
		m.UpdateViewport()
		return m, waitForToggle(m.toggleCh)
	case newEvent:
		if msg.Event.IsStart() {
			m.data.clean()
			if m.serverMode {
				m.data.onlyStepsMode = true
			}
			m.running = true
			m.UpdateViewport()
			return m, readNextEvent(m.runner)
		}
		if !msg.WaitingMore {
			m.running = false
			if m.serverMode {
				return m, readNextEvent(m.runner)
			}
			return m, nil
		}
		m.data.add(msg.Event)
		m.UpdateViewport()
		return m, readNextEvent(m.runner)
	}

	return m, cmd
}

func (m *model) UpdateViewport() {
	msg := m.window

	m.table.Headers(m.data.headers()...)
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
	if m.serverMode && len(m.data.events) == 0 {
		return "Waiting for connection on " + m.db + "...\n\nPress q to quit."
	}

	menuBindings := theme.DefaultKB.Menu()
	if m.serverMode {
		menuBindings = theme.DefaultKB.ServerMenu()
	}
	menu := "  " + m.help.ShortHelpView(menuBindings)
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
