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
	//"golang.org/x/term"

	"github.com/finitelycraig/nethackathon-charity-progress/db"
)

type Progress struct {
	fundraiser db.Fundraiser
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
	maxWidth = 80 //,_,_ = term.GetSize(int(os.Stdout.Fd()))
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func Run() {
	p := NewProgress()

	if _, err := tea.NewProgram(p, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

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
		//cmd := m.progress.IncrPercent(0.25)
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

func (p Progress) amountString(pad string) string {
    raised,err := strconv.Atoi(p.fundraiser.Raised)
    if err != nil {
        return "\n\n"
    }
    supporters := p.fundraiser.Supporters
    supporterors := "supporters"
    if supporters == 1 {
        supporterors = "supporter"
    }
    return fmt.Sprintf("%s$%.2f USD\n%sRaised by %d %s\n", pad, float64(raised)/100.00, pad, supporters, supporterors)
}

func (p Progress) goalString(pad string) string {
    goal,err := strconv.Atoi(p.fundraiser.GoalAmount)
    if err != nil {
        return "\n\n"
    }
    return fmt.Sprintf("$%.2f USD%s\nGoal%s\n", float64(goal)/100.00, pad, pad)
}

func (m Progress) View() string {
	pad := strings.Repeat(" ", padding)
    raisedSummaryStyle := lipgloss.NewStyle().Width(m.progress.Width/2).Align(lipgloss.Left)
    goalSummaryStyle := lipgloss.NewStyle().Width(m.progress.Width/2).Align(lipgloss.Right)
    progressSummary := lipgloss.JoinHorizontal(lipgloss.Center, raisedSummaryStyle.Render(m.amountString(pad)), goalSummaryStyle.Render(m.goalString(pad+" ")))
    return "\n" +
        progressSummary + "\n" +  
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press ctl+c or q or Esc to quit")
}

func (p Progress) tickCmd() tea.Cmd {
    p.fundraiser = db.GetFundraiserData()
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
