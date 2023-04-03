package ui

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gookit/slog"
	"github.com/robbailey3/openai-cli/openai"
)

type (
	errMsg error
)

type Model struct {
	messages     []openai.ChatMessage
	viewport     viewport.Model
	textarea     textarea.Model
	senderStyle  lipgloss.Style
	err          error
	openAiClient openai.Client
}

func NewModel() Model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)
	return Model{
		textarea:     ta,
		messages:     []openai.ChatMessage{},
		viewport:     vp,
		senderStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		openAiClient: openai.NewClient(),
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.SendMessage(m.textarea.Value())
			msgStr := ""
			for _, mess := range m.messages {
				msgStr += mess.Content + "\n"
			}
			m.viewport.SetContent(msgStr)
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func (m Model) SendMessage(message string) {
	m.messages = append(m.messages, openai.ChatMessage{
		Role:    "user",
		Content: message,
	})

	completion, err := m.openAiClient.GetChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:            "gpt-4",
		Messages:         m.messages,
		Temperature:      1,
		TopP:             1,
		N:                1,
		Stream:           false,
		MaxTokens:        250,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
	})
	if err != nil {
		slog.Error(err)
		return
	}

	m.messages = append(m.messages, completion.Choices[0].Message)
}
