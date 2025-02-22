package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title, description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func initialModel() model {
	items := []list.Item{
		item{title: "git status", description: "show working tree status"},
		item{title: "git add .", description: "add all changes to staging"},
		item{title: "git commit", description: "commit staged changes"},
		item{title: "git push", description: "push commits to remote"},
		item{title: "git pull", description: "pull changes from remote"},
		item{title: "git branch", description: "list branches"},
		item{title: "git checkout", description: "switch branches"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "april"
	l.SetShowHelp(false)

	return model{
		list: l,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func executeGitCommand(command string) error {
	cmd := exec.Command("git", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "ctrl+n":
			m.list.CursorDown()
		case "ctrl+p":
			m.list.CursorUp()
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.title
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return "quit april\n"
	}
	return docStyle.Render(m.list.View())
}

func main() {
	// fmt.Println("april")
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	if model, ok := m.(model); ok && model.choice != "" {
		fmt.Printf("april: %s\n", model.choice)
		executeGitCommand(model.choice)
	}
}

