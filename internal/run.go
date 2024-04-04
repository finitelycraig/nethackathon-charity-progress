package internal

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/finitelycraig/nethackathon-charity-progress/data"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/muesli/termenv"
)

func Run() {
	p := data.NewProgress()
	if _, err := tea.NewProgram(p, tea.WithAltScreen()).Run(); err != nil {
		log.Error().Err(err).Msg("Could not start local progress bar")
		os.Exit(1)
	}
}

const (
	host = "localhost"
	port = "42069"
)

type app struct {
	*ssh.Server
}

func newApp() *app {
	a := new(app)
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.MiddlewareWithProgramHandler(a.ProgramHandler, termenv.ANSI256),
			activeterm.Middleware(),
		),
	)
	if err != nil {
		log.Error().Err(err).Msg("Could not start new app")
	}

	a.Server = s
	return a
}

func (a *app) Start() {
	var err error
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info().Msg(fmt.Sprintf("Starting SSH server host %s port %s", host, port))
	go func() {
		if err = a.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("Could not start server in app.Start()")
			done <- nil
		}
	}()

	<-done
	log.Info().Msg("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := a.Shutdown(ctx); err != nil {
		log.Error().Msg("Could not stop server")
	}
}

func (a *app) ProgramHandler(s ssh.Session) *tea.Program {
	log.Info().Msg(fmt.Sprintf("Starting a room without a host"))
	p := data.NewProgress()
	programOptions := bubbletea.MakeOptions(s)
	programOptions = append(programOptions, tea.WithAltScreen())
	prog := tea.NewProgram(p, programOptions...)

	return prog
}

func RunOverSSH() {
	app := newApp()
	app.Start()
}
