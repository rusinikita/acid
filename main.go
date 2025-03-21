package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/db"
	"github.com/rusinikita/acid/runner"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/viewport"
)

/*
design

---------------------------------------------
	[x]	step n			| rows affected - 0
	[x] step n			| -------------------------------------
	[\] step n [tx2]	| [tx 1]
	[x] step n [tx1]-	| Select * from bla bla
	step n				| where biba = boba
	step n				| --------
	step n				| result
	step n				|     col 1 | col 2 |
	step n				| 0 | val 1 | val 2 |
	step n				| 1 | val 1 | val 2 |
	step n				| -------------------------------------
	step n				| [tx 2]
	step n				| Select * from bla bla
	step n				| --------
	step n				| [/] Running response
	step n				|
	step n				|
---------------------------------------------
enter - next step

плохо - нужно возвращаться к ожидающему запросу чтобы посмотреть результат - решение - табы-транзакции (общий, tx1, tx2...)
*/

/*
  common  	|  - tx1 -  	| tx2
	step	|				|
			----------------
			|	begin		|
			|	step 2		|
			|				| step 4
			|	step 3		|	- waiting
			|	commit		| 	- waiting
			----------------	- waiting
			|				| step 4 result
			|				|
------------				------------
|
|
|

Row - main
Transaction events in row
Event form: short/detail
Event type: begin/commit/deadlock/select/select_result/insert_result
Event tx selected: common, tx1, tx2

Short form:
step 4 - select name
select * from...

Detail form:

*/

func main() {
	tx1 := call.TrxID("1")
	tx2 := call.TrxID("2")

	sequence := call.Sequence{Calls: []call.Step{
		call.Call("CREATE TABLE IF NOT EXISTS exec_test (id SERIAL PRIMARY KEY, name TEXT)"),
		call.Begin(tx1),
		call.Begin(tx2),
		call.Call("insert into exec_test (name) values ('biba')", tx1),
		call.Call("select * from exec_test", tx2),
		call.Commit(tx1),
		call.Call("select * from exec_test", tx2),
		call.Call("select * from exec_test"),
	}}

	next := make(chan struct{})

	go func() {
		for {
			fmt.Scanln()
			next <- struct{}{}
		}
	}()

	c := runner.New(db.Connect()).Run(sequence)
	for event := range c {
		fmt.Println(event.View())
	}

	//p := tea.NewProgram(
	//	model{},
	//	tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
	//	tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	//)
	//
	//_, err := p.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
}

type model struct {
	tick     int
	ready    bool
	viewport viewport.Model
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return true
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if msg, ok := msg.(bool); ok && msg {
		m.tick++
		m.viewport.SetContent(m.content())
		cmds = append(cmds, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return true
		}))
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.content())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) content() string {
	sb := strings.Builder{}

	for i := 0; i < m.tick*10; i++ {
		sb.WriteString(strconv.Itoa(i) + "\n")
	}

	return sb.String()
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

func (m model) headerView() string {
	title := titleStyle.Render("Mr. Pager")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
