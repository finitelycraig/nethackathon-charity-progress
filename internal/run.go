package internal

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitelycraig/nethackathon-charity-progress/data"
	"os"
)

func Run() {
	p := data.NewProgress()
	if _, err := tea.NewProgram(p, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}
