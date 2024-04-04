package data

import (
	"fmt"
	"os"
    "strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

    "github.com/finitelycraig/nethackathon-charity-progress/db"
)

type Progress struct {
    fundraiser db.Fundraiser
	progress   progress.Model
}

func (p Progress) goalPercentage() float64 {
    current,err := float64(strconv.Atoi(fundraiser.Raised))
    if err != nil {
        return 0.0
    }
    goal, err := float64(strconv.Atoi(fundraiser.Raised))
    if err != nil {
        return 0.0
    }
    _ :=
    return 0.5 
}

func NewProgress() Progress {
    return Progress{fundraiser: db.GetFundraiserData(), progress: progress.New(progress.WithDefaultGradient())}
}

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func Run() {
	p := NewProgress()

	if _, err := tea.NewProgram(p).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

func (p Progress) Init() tea.Cmd {
	return tickCmd()
}

func (m Progress) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(0.25)
        cmd = progress.Model.SetPercent(p.goalPercentage())
		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m Progress) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*30, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
