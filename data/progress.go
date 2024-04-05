package data

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/finitelycraig/nethackathon-charity-progress/db"
)

type Progress struct {
	fundraiser db.Fundraiser
	width      int
	height     int
	progress   progress.Model
}

func (p Progress) goalPercentage() float64 {
	current, err := strconv.Atoi(p.fundraiser.Raised)
	if err != nil {
		return 0.0
	}
	var goal int
	goal, err = strconv.Atoi(p.fundraiser.GoalAmount)
	if err != nil {
		return 0.0
	}
	return (float64(current) / float64(goal))
}

func NewProgress() Progress {
	return Progress{fundraiser: db.GetFundraiserData(), progress: progress.New(progress.WithDefaultGradient())}
}

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type tickMsg time.Time

func (p Progress) Init() tea.Cmd {
	return p.tickCmd()
}

func (m Progress) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "q":
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
        m.fundraiser = db.GetFundraiserData()        
		cmd := m.progress.SetPercent(m.goalPercentage())
		return m, tea.Batch(m.tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
	return m, nil
}

func (p Progress) amountString() string {
	raised, err := strconv.Atoi(p.fundraiser.Raised)
	if err != nil {
		return "\n\n"
	}
	supporters := p.fundraiser.Supporters
	supporterors := "supporters"
	if supporters == 1 {
		supporterors = "supporter"
	}
	return fmt.Sprintf("$%.2f USD\nRaised by %d %s\n", float64(raised)/100.00, supporters, supporterors)
}

func (p Progress) goalString() string {
	goal, err := strconv.Atoi(p.fundraiser.GoalAmount)
	if err != nil {
		return "\n\n"
	}
	return fmt.Sprintf("$%.2f USD\nGoal\n", float64(goal)/100.00)
}

func (m Progress) View() string {
	raisedSummaryStyle := lipgloss.NewStyle().Width(m.progress.Width / 2).Align(lipgloss.Left)
	goalSummaryStyle := lipgloss.NewStyle().Width(m.progress.Width / 2).Align(lipgloss.Right)
	progressSummary := lipgloss.JoinHorizontal(lipgloss.Center, raisedSummaryStyle.Render(m.amountString()), goalSummaryStyle.Render(m.goalString()))
	progressBar := m.progress.View() //+ "\n\n" + helpStyle("Press ctl+c or q or Esc to quit")
	view := lipgloss.JoinVertical(lipgloss.Center, progressSummary, progressBar)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
}

func (p Progress) tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
