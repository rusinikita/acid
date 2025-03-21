package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/event"
	"slices"
)

type model struct {
	runner runner

	running bool
	window  tea.WindowSizeMsg

	table *table.Table
	vp    viewport.Model
	data  *eventData
}

func NewRunTable(runner runner) tea.Model {
	data := &eventData{}

	t := table.New().Data(data)

	m := &model{
		runner:  runner,
		running: true,

		table: t,
		data:  data,
	}

	styleFunc := func(row, _ int) lipgloss.Style {
		switch {
		case row == table.HeaderRow:
			return headerStyle
		default:
			return oddRowStyle.Width(m.window.Width/m.data.Columns() - m.data.Columns())
		}
	}
	t.StyleFunc(styleFunc)

	return m
}

func (m *model) Init() tea.Cmd {
	m.runner.Run()

	return readNextEvent(m.runner)
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
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.table.Height(msg.Height)
		m.table.Width(msg.Width)
		m.vp.Width = msg.Width
		m.vp.Height = msg.Height

		m.window = msg

		m.UpdateViewport()

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

// windowWidth - dividerWidth * (columns - 1)
// 10% width per not selected column and the rest for selected

func (m *model) View() string {
	return m.vp.View()
}

type newEvent struct {
	Event       event.Event
	WaitingMore bool
}

type eventData struct {
	events             []event.Event
	transactions       []call.TrxID
	selectedTrxIndex   call.TrxID
	selectedEventIndex int
}

func (ed *eventData) add(e event.Event) {
	ed.events = append(ed.events, e)

	if !slices.Contains(ed.transactions, e.Trx()) {
		ed.transactions = append(ed.transactions, e.Trx())
	}
}

func (ed *eventData) headers() []string {
	headers := make([]string, len(ed.transactions))
	for i, trx := range ed.transactions {
		if trx == "" {
			headers[i] = "No tx requests"
			continue
		}

		headers[i] = "TX - " + string(trx)
	}

	return headers
}

func (ed *eventData) At(row, cell int) string {
	e := ed.events[row]
	trx := ed.transactions[cell]

	s := e.Cell(trx)

	return s
}

func (ed *eventData) Rows() int {
	return len(ed.events)
}

func (ed *eventData) Columns() int {
	return len(ed.transactions)
}
