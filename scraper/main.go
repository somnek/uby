package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	page    int
	repo    string
	input   textinput.Model
	deps    []Dep
	err     error
	spinner spinner.Model
	done    bool
	count   int
}

type (
	errMsg error
)

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter a dependency graph URL..."
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 200

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		input:   ti,
		page:    0,
		err:     nil,
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmdSpinner tea.Cmd
	var cmdInput tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			if m.page == 0 {
				// todo: validate/clean input
				m.page = 1
				m.repo = m.input.Value()
				// hardcodeUrl := "https://github.com/aquasecurity/trivy/network/dependents"
				return m, tea.Batch(InitScrape(), m.spinner.Tick)
			}

		}

	case errMsg:
		m.err = msg
		return m, nil

	case Done:
		m.done = true
		m.count++
		return m, nil

	case InitScrapeTick:
		return m, nil

	case scrapeTick:
		nextUrl, nextPageExist := msg
		return m, Batch(Scrape(nextUrl), m.spinner.Tick)
	}

	m.spinner, cmdSpinner = m.spinner.Update(msg)
	m.input, cmdInput = m.input.Update(msg)

	return m, tea.Batch(cmd, cmdSpinner, cmdInput)

}

func (m model) View() string {
	title := "🥦 Uby"
	body := ""
	footer := "(esc to quit)"

	switch m.page {
	case 0:
		body = m.input.View()
	case 1:
		if !m.done {
			body += m.spinner.View()
			body += "Fetching dependencies..."
		} else {
			body = "Done!"
		}
	}

	outText := fmt.Sprintf("%s\n\n%s\n\n%s\n%d", title, body, footer, m.count)
	return outText

}
