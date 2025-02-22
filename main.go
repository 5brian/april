package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
		item{title: "git commit -m", description: "commit staged changes"},
		item{title: "git push", description: "push commits to remote"},
		item{title: "git pull", description: "pull changes from remote"},
		item{title: "git branch", description: "list branches"},
		item{title: "git checkout -b", description: "create and switch to new branch"},
		item{title: "git checkout main", description: "switch to main branch"},
		item{title: "git log", description: "show commit logs"},
		item{title: "git fetch", description: "download objects and refs from remote"},
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
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	if parts[1] == "commit" && parts[2] == "-m" {
		cmd := exec.Command("git", "commit", "-m", strings.Join(parts[3:], " "))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	cmd := exec.Command("git", parts[1:]...)
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
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	if model, ok := m.(model); ok && model.choice != "" {
		fmt.Printf("april: %s\n", model.choice)
		if strings.Contains(model.choice, "commit -m") {
			fmt.Print("Enter commit message: ")
			reader := bufio.NewReader(os.Stdin)
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading commit message: %v\n", err)
				return
			}
			message = strings.TrimSpace(message)
			err = executeGitCommand(fmt.Sprintf("git commit -m %s", message))
		} else if strings.Contains(model.choice, "checkout -b") {
			fmt.Print("Enter branch name: ")
			reader := bufio.NewReader(os.Stdin)
			branchName, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading branch name: %v\n", err)
				return
			}
			branchName = strings.TrimSpace(branchName)
			err = executeGitCommand(fmt.Sprintf("git checkout -b %s", branchName))
		} else {
			err = executeGitCommand(model.choice)
		}

		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		}
	}
}
