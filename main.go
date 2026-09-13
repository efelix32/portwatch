// Command portwatch shows a live, terminal-based view of the local
// machine's network connections.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/efelix32/portwatch/internal/ui"
)

func main() {
	p := tea.NewProgram(ui.New())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: %v\n", err)
		os.Exit(1)
	}
}
