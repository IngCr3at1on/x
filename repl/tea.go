package repl

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	model struct {
		viewport  viewport.Model
		textinput textinput.Model

		repl *Repl
	}
)

var _ tea.Model = new(model)

func initialModel(r *Repl) model {
	m := model{
		repl:      r,
		textinput: textinput.New(),
	}

	m.textinput.Prompt = r.prompt
	m.textinput.Focus()

	return m
}

func (model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.repl.isQuitMessage(message) {
		return m, tea.Quit
	}

	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			cmd, err := m.repl.eval(m.repl.ctx, m.textinput.Value(), m.repl)
			if err != nil {
				m.error(err)
			}
			if cmd != nil {
				cmds = append(cmds, cmd)
			}

			m.textinput.Reset()

			m.viewport.SetContent(m.repl.buf.String())
			m.viewport.GotoBottom()

			if m.repl.highPerf {
				cmds = append(cmds, viewport.Sync(m.viewport))
			}
		}

	case tea.WindowSizeMsg:
		if !m.repl.ready {
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.HighPerformanceRendering = m.repl.highPerf

			m.viewport.SetContent(m.repl.buf.String())
			m.viewport.GotoBottom()

			m.repl.ready = true
		}

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height

		m.textinput.Width = msg.Width

		if m.repl.highPerf {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(message)
	cmds = append(cmds, cmd)

	if m.repl.ready {
		m.textinput, cmd = m.textinput.Update(message)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.repl.ready {
		return "\n Initializing..."
	}

	return lipgloss.JoinVertical(lipgloss.Top, m.viewport.View(), m.textinput.View())
}

func (m model) error(err error) {
	if err != nil {
		if m.repl.ready {
			_, _err := m.repl.buf.WriteString(err.Error())
			if _err != nil {
				m.repl.error(_err)
				return
			}
		} else {
			m.repl.error(err)
		}
	}
}
