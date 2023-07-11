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
	logs    string
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
		count:   0,
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
				// <-- ️ HARD CODED URL FOR TESTING ⚠️
				hardcodeUrl := "https://github.com/aquasecurity/trivy/network/dependents"
				m.repo = hardcodeUrl
				// ⚠️ HARD CODED URL FOR TESTING -->
				return m, tea.Batch(InitScrape(), m.spinner.Tick)
			}

		}

	case errMsg:
		m.err = msg
		return m, nil

	case InitScrapeTick:
		return m, tea.Batch(Scrape(m.repo), m.spinner.Tick)

	case ScrapeTick:
		nextUrl := msg.nextUrl
		if nextUrl != "" {
			m.count++
			m.logs += fmt.Sprintf("📦 %s\n", nextUrl)
			return m, tea.Batch(Scrape(nextUrl), m.spinner.Tick)
		} else {
			m.logs += "🧨"
			return m, tea.Quit
		}
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
			body += fmt.Sprintf("Fetching dependencies from %s...", m.repo)
		} else {
			body = "Done!"
		}
	}

	outText := fmt.Sprintf("%s\n\n%s\n\n%s\ncount: %d\n\n%s\n\n", title, body, footer, m.count, m.logs)
	return outText

}
