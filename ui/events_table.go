package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/event"
)

type model struct {
	runner runner

	events             []event.Event
	running            bool
	selectedTrx        call.TrxID
	selectedEventIndex int
}

func newRunTable(runner runner) model {
	return model{
		runner:  runner,
		running: true,
	}
}

func (m model) Init() tea.Cmd {
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
	case tea.WindowSizeMsg:
	case newEvent:
		m.running = msg.WaitingMore
		if !m.running {
			return m, nil
		}

		m.events = append(m.events, msg.Event)
		return m, readNextEvent(m.runner)
	}

	return m, nil
}

func (m model) View() string {
	//TODO implement me
	panic("implement me")
}

type newEvent struct {
	Event       event.Event
	WaitingMore bool
}
