package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	url    string
	status int
	err    error
}

func (m model) checkServer() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(m.url)

	if err != nil {
		return errMsg{err}
	}

	return statusMsg(res.StatusCode)
}

type statusMsg int

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func (m model) Init() tea.Cmd {
	return m.checkServer
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case statusMsg:
		m.status = int(msg)
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case tea.KeyMsg:
		// catches ctrl c and the like
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	}

	// if we get other messages, don't do anything
	return m, nil
}

func (m model) View() string {
	// if error, print
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	s := fmt.Sprintf("Checking %s ... ", m.url)

	// when the server responds with a status, add it to line
	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}

	return "\n" + s + "\n\n"
}

func main() {
	var initialModel model

	if len(os.Args) > 1 {
		initialModel.url = os.Args[1]
	} else {
		initialModel.url = "https://charm.sh/"
	}

	if _, err := tea.NewProgram(initialModel).Run(); err != nil {
		fmt.Printf("Uh oh, there as an error: %v\n", err)
		os.Exit(1)
	}
}
