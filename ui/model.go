package ui

import (
	"fmt"
	"strings"
	"time"

	dockerwrapper "github.com/ahmedmaaloul/godock-tui-manager/docker"
	"github.com/ahmedmaaloul/godock-tui-manager/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// Define messages
type errMsg error
type containersMsg []dockerwrapper.Container

// Model for Bubble Tea
type Model struct {
	dockerClient *dockerwrapper.Client
	containers   []dockerwrapper.Container
	cursor       int
	err          error
	loading      bool
	message      string // status message at the bottom
}

// NewModel initializes the UI model
func NewModel(dc *dockerwrapper.Client) Model {
	return Model{
		dockerClient: dc,
		containers:   []dockerwrapper.Container{},
		cursor:       0,
		loading:      true,
	}
}

// Init run on startup
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchContainers(),
		tea.EnterAltScreen,
	)
}

// Cmd to fetch containers
func (m Model) fetchContainers() tea.Cmd {
	return func() tea.Msg {
		if m.dockerClient == nil {
			return errMsg(fmt.Errorf("docker client not initialized"))
		}
		containers, err := m.dockerClient.ListContainers()
		if err != nil {
			return errMsg(err)
		}
		return containersMsg(containers)
	}
}

// Cmd to start a container
func (m Model) startContainer(id string) tea.Cmd {
	return func() tea.Msg {
		err := m.dockerClient.StartContainer(id)
		if err != nil {
			return errMsg(err)
		}
		// Refresh list after action
		return m.fetchContainers()()
	}
}

// Cmd to stop a container
func (m Model) stopContainer(id string) tea.Cmd {
	return func() tea.Msg {
		err := m.dockerClient.StopContainer(id)
		if err != nil {
			return errMsg(err)
		}
		return m.fetchContainers()()
	}
}

// Update loop
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else if len(m.containers) > 0 {
				m.cursor = len(m.containers) - 1 // wrap around
			}

		case "down", "j":
			if m.cursor < len(m.containers)-1 {
				m.cursor++
			} else {
				m.cursor = 0 // wrap around
			}

		case "r":
			m.loading = true
			m.message = "Refreshing..."
			return m, m.fetchContainers()

		case "s":
			if len(m.containers) > 0 {
				id := m.containers[m.cursor].ID
				m.message = fmt.Sprintf("Starting %s...", id)
				m.loading = true
				return m, m.startContainer(id)
			}

		case "x":
			if len(m.containers) > 0 {
				id := m.containers[m.cursor].ID
				m.message = fmt.Sprintf("Stopping %s...", id)
				m.loading = true
				return m, m.stopContainer(id)
			}
		}

	case containersMsg:
		m.containers = msg
		m.loading = false
		m.message = fmt.Sprintf("Updated at %s", time.Now().Format("15:04:05"))
		// Adjust cursor if out of bounds (e.g. list shrank)
		if m.cursor >= len(m.containers) && len(m.containers) > 0 {
			m.cursor = len(m.containers) - 1
		}

	case errMsg:
		m.err = msg
		m.loading = false
		m.message = fmt.Sprintf("Error: %v", msg)
	}

	return m, nil
}

// View container table
func (m Model) View() string {
	s := strings.Builder{}

	s.WriteString("🐳 GoDock Manager\n\n")

	// Header
	s.WriteString(styles.HeaderStyle.Render(fmt.Sprintf("%-12s %-30s %-20s %-15s", "ID", "IMAGE", "STATUS", "STATE")) + "\n")

	// Rows
	for i, c := range m.containers {
		cursor := " " // no cursor
		rowStyle := styles.RowStyle

		if m.cursor == i {
			cursor = ">" // cursor
			rowStyle = styles.SelectedRowStyle
		}

		// Truncate Image name if too long
		imageName := c.Image
		if len(imageName) > 30 {
			imageName = imageName[:27] + "..."
		}

		// Colorize state
		stateStr := c.State
		if c.State == "running" {
			stateStr = styles.StatusRunning.Render(c.State)
		} else if c.State == "exited" {
			stateStr = styles.StatusExited.Render(c.State)
		} else {
			stateStr = styles.StatusPaused.Render(c.State) // paused or other
		}

		// Render the row
		rowStr := fmt.Sprintf("%s %-12s %-30s %-20s %s", cursor, c.ID, imageName, c.Status, stateStr)
		s.WriteString(rowStyle.Render(rowStr) + "\n")
	}

	if len(m.containers) == 0 {
		if m.loading {
			s.WriteString("\n Loading containers...\n")
		} else {
			s.WriteString("\n No containers found.\n")
		}
	}

	// Footer / Status
	s.WriteString("\n")
	if m.message != "" {
		s.WriteString(m.message + "\n")
	}

	help := "↑/k: Up • ↓/j: Down • s: Start • x: Stop • r: Refresh • q: Quit"
	s.WriteString(styles.HelpStyle.Render(help))

	return styles.AppStyle.Render(s.String())
}
