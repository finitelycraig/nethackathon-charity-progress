package data

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
	"github.com/finitelycraig/nethackathon-charity-progress/db"
	"github.com/skip2/go-qrcode"
)

type Progress struct {
	fundraiser db.Fundraiser
	config	   viper.Viper
	width      int
	height     int
	progress   progress.Model
	qr	   string
	showInfo   bool
	showLink   bool
	showBitly  bool
	showQR     bool
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
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if _, ok := viper.Get("tagline").(string); !ok {
		panic(fmt.Errorf("tagline not set as string in config.yaml"))
	} else if _, ok = viper.Get("link").(string); !ok {
		panic(fmt.Errorf("link not set as string in config.yaml"))
	}
	return Progress{fundraiser: db.GetFundraiserData(), progress: progress.New(progress.WithDefaultGradient()), qr: generateQR(viper.Get("link").(string))}
}

const (
	padding  = 2
	maxWidth = 80
)

var infoStyle = lipgloss.NewStyle().Align(lipgloss.Center).Render

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
			case "i":
				m.showInfo = !m.showInfo
			case "l":
				m.showLink = !m.showLink
			case "b":
				m.showBitly = !m.showBitly
			case "c":
				m.showQR = !m.showQR
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

func quarterBlockFromBitmap(upperLeft, upperRight, lowerLeft, lowerRight bool) string {
	if upperLeft && upperRight && lowerLeft && lowerRight {
		return "█"
	}
	if upperLeft && upperRight && !lowerLeft && !lowerRight {
		return "▀"
	}
	if !upperLeft && !upperRight && lowerLeft && lowerRight {
		return "▄"
	}
	if upperLeft && !upperRight && lowerLeft && !lowerRight {
		return "▌"
	}
	if !upperLeft && upperRight && !lowerLeft && lowerRight {
		return "▐"
	}
	if upperLeft && !upperRight && !lowerLeft && lowerRight {
		return "▚"
	}
	if !upperLeft && upperRight && lowerLeft && !lowerRight {
		return "▞"
	}
	if upperLeft && upperRight && lowerLeft && !lowerRight {
		return "▛"
	}
	if upperLeft && upperRight && !lowerLeft && lowerRight {
		return "▜"
	}
	if upperLeft && !upperRight && lowerLeft && lowerRight {
		return "▙"
	}
	if !upperLeft && upperRight && lowerLeft && lowerRight {
		return "▟"
	}
	if upperLeft && !upperRight && !lowerLeft && !lowerRight {
		return "▘"
	}
	if !upperLeft && upperRight && !lowerLeft && !lowerRight {
		return "▝"
	}
	if !upperLeft && !upperRight && lowerLeft && !lowerRight {
		return "▖"
	}
	if !upperLeft && !upperRight && !lowerLeft && lowerRight {
		return "▗"
	}
	if !upperLeft && !upperRight && !lowerLeft && !lowerRight {
		return " "
	}
	return ""
}

func generateQR(link string) string{
	qr, err := qrcode.New(link, qrcode.Medium)
	if err != nil {
		return "Unable to generate QR code"
	}
	var output string
	for i := 3; i < len(qr.Bitmap()) - 3; i+=2 {
		for j := 3; j < len(qr.Bitmap()) -3; j+=2 {
			upperLeft := qr.Bitmap()[i][j]
			upperRight := qr.Bitmap()[i][j+1]
			lowerLeft := qr.Bitmap()[i+1][j]
			lowerRight := qr.Bitmap()[i+1][j+1]
			output += quarterBlockFromBitmap(upperLeft, upperRight, lowerLeft, lowerRight)
		}
		output += "\n"
	}
	return output
}

func (m Progress) View() string {
	raisedSummaryStyle := lipgloss.NewStyle().Width(m.progress.Width / 2).Align(lipgloss.Left)
	goalSummaryStyle := lipgloss.NewStyle().Width(m.progress.Width / 2).Align(lipgloss.Right)
	progressSummary := lipgloss.JoinHorizontal(lipgloss.Center, raisedSummaryStyle.Render(m.amountString()), goalSummaryStyle.Render(m.goalString()))
	progressBar := m.progress.View() //+ "\n\n" + infoStyle("Press ctl+c or q or Esc to quit")
	var view string
	if m.showInfo {
		if tagline, ok := viper.Get("tagline").(string); ok {
			progressBar += "\n\n" + infoStyle(tagline)
		}
		if m.showLink && !m.showBitly {
			if link, ok := viper.Get("link").(string); ok {
				progressBar += "\n" + infoStyle(link)
			}
		} else if m.showLink && m.showBitly {
			if shortlink, ok := viper.Get("shortlink").(string); ok {
				progressBar += "\n" + infoStyle(shortlink)
			}
		}
	}
	if m.showQR {
		view = infoStyle(m.qr)
	} else {
		view = lipgloss.JoinVertical(lipgloss.Center, progressSummary, progressBar)
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
}

func (p Progress) tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
