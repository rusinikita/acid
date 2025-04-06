package list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rusinikita/acid/sequence"
	"github.com/rusinikita/acid/ui/router"
	"io"
	"strings"
)

type choose struct {
	err  error
	list list.Model
}

func New(db string) tea.Model {
	listModel := list.New(nil, itemDelegate{}, 20, 14)
	listModel.SetFilteringEnabled(false)
	listModel.SetShowStatusBar(false)
	listModel.Title = fmt.Sprintf("Select sequence to run on '%s'", db)
	listModel.InfiniteScrolling = true

	items := make([]list.Item, 0, len(sequence.Sequences))
	for _, task := range sequence.Sequences {
		items = append(items, item(task))
	}
	listModel.SetItems(items)

	return choose{
		list: listModel,
	}
}

func (c choose) Init() tea.Cmd {
	return nil
}

func (c choose) Update(msg tea.Msg) (m tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		c.list.SetSize(
			msg.Width-listStyle.GetHorizontalFrameSize(),
			msg.Height-listStyle.GetVerticalFrameSize(),
		)

	case tea.KeyMsg:
		if msg.String() == "enter" {
			return m, router.Route("run", c.list.Index())
		}
	}

	c.list, cmd = c.list.Update(msg)

	return c, cmd
}

func (c choose) View() string {
	if c.err != nil {
		return c.err.Error()
	}

	return listStyle.Render(c.list.View())
}

type item sequence.Sequence

func (i item) FilterValue() string {
	return i.Name
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 6 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Name)

	fn := titleStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedTitleStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(str), "\n", descStyle.Render(i.Description))
}

var (
	listStyle          = lipgloss.NewStyle().Margin(1, 2)
	descStyle          = lipgloss.NewStyle().Padding(1, 4)
	titleStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedTitleStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)
