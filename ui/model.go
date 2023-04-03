package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/robbailey3/openai-cli/openai"
)

type Model struct {
	Messages []openai.ChatMessage
}

func NewModel() Model {
	return Model{
		Messages: []openai.ChatMessage{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	// TODO implement me
	panic("implement me")
}
