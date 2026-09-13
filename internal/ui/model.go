// Package ui renders the portwatch connection table as a bubbletea program.
package ui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/efelix32/portwatch/internal/scanner"
)

// RefreshInterval controls how often the connection table is rescanned.
const RefreshInterval = 2 * time.Second

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	footerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type tickMsg time.Time
type scanResultMsg []scanner.Connection
type scanErrMsg error

// Model is the bubbletea model for the connection table.
type Model struct {
	table table.Model
	err   error
}

// New builds the initial model.
func New() Model {
	columns := []table.Column{
		{Title: "PROTO", Width: 5},
		{Title: "LOCAL", Width: 22},
		{Title: "REMOTE", Width: 22},
		{Title: "STATUS", Width: 12},
		{Title: "PID", Width: 8},
		{Title: "PROCESS", Width: 18},
	}
	t := table.New(table.WithColumns(columns), table.WithFocused(true), table.WithHeight(20))
	return Model{table: t}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(scanCmd(), tickCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	case tickMsg:
		return m, tea.Batch(scanCmd(), tickCmd())
	case scanResultMsg:
		m.err = nil
		m.table.SetRows(rowsFromConnections(msg))
	case scanErrMsg:
		m.err = msg
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	view := headerStyle.Render("portwatch") + footerStyle.Render(" -- live local connections\n\n")
	view += m.table.View() + "\n"
	if m.err != nil {
		view += errorStyle.Render(fmt.Sprintf("\nerror: %v\n", m.err))
	}
	view += footerStyle.Render("\nq to quit -- refreshes every " + RefreshInterval.String())
	return view
}

func rowsFromConnections(conns []scanner.Connection) []table.Row {
	rows := make([]table.Row, 0, len(conns))
	for _, c := range conns {
		rows = append(rows, table.Row{
			c.Proto,
			c.LocalAddr,
			c.RemoteAddr,
			c.Status,
			strconv.Itoa(int(c.PID)),
			c.ProcessName,
		})
	}
	return rows
}

func tickCmd() tea.Cmd {
	return tea.Tick(RefreshInterval, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func scanCmd() tea.Cmd {
	return func() tea.Msg {
		conns, err := scanner.Scan()
		if err != nil {
			return scanErrMsg(err)
		}
		return scanResultMsg(conns)
	}
}
