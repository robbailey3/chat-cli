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
	errMsg        error
	completionMsg openai.ChatMessage
)

type Model struct {
	messages     []openai.ChatMessage
	viewport     viewport.Model
	textarea     textarea.Model
	senderStyle  lipgloss.Style
	err          error
	openAiClient openai.Client
	ready        bool
	isLoading    bool
}

func NewModel() Model {
	ta := textarea.New()
	ta.Placeholder = ""
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	ta.KeyMap.InsertNewline.SetEnabled(false)
	return Model{
		textarea:     ta,
		messages:     []openai.ChatMessage{},
		senderStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		openAiClient: openai.NewClient(),
		ready:        false,
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
	case tea.WindowSizeMsg:
		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-m.textarea.Height()-2)
			m.textarea.SetWidth(msg.Width)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - m.textarea.Height() - 2
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			if m.isLoading {
				return m, nil
			}
			m.isLoading = true
			m.messages = append(m.messages, openai.ChatMessage{
				Content: m.textarea.Value(),
				Role:    "user",
			})
			return m, m.SendMessage(m.textarea.Value())
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	case completionMsg:
		m.isLoading = false
		msgStr := ""
		m.messages = append(m.messages, openai.ChatMessage(msg))
		for _, mess := range m.messages {
			msgStr += m.getStyledMessage(mess) + "\n"
		}
		m.viewport.SetContent(msgStr)
		m.textarea.Reset()
		m.viewport.GotoBottom()
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) getStyledMessage(msg openai.ChatMessage) string {
	if msg.Role == "user" {
		row := lipgloss.NewStyle().Width(m.viewport.Width).AlignHorizontal(lipgloss.Right).PaddingTop(1).PaddingRight(2)
		chatBox := lipgloss.NewStyle().Width(m.viewport.Width / 2).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#04B575"))

		return row.Render(chatBox.Render(msg.Content))
	}
	row := lipgloss.NewStyle().Width(m.viewport.Width).AlignHorizontal(lipgloss.Left).PaddingLeft(2)
	chatBox := lipgloss.NewStyle().Width(m.viewport.Width / 2).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#11349c"))
	return row.Render(chatBox.Render(msg.Content))
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	)
}

func (m Model) SendMessage(message string) tea.Cmd {
	return func() tea.Msg {
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
			return nil
		}
		if len(completion.Choices) > 0 {
			return completionMsg(completion.Choices[0].Message)
		}

		return nil
	}
}
