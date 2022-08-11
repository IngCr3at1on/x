package repl

import tea "github.com/charmbracelet/bubbletea"

const defaultPrompt = "$ "

func defaultIsQuitMessageF(message tea.Msg) bool {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return true
		case tea.KeyCtrlC:
			return true
		}
	}

	return false
}

type Option func(*Repl)

func WithPrompt(prompt string) Option {
	return func(r *Repl) {
		r.prompt = prompt
	}
}

func WithIsQuitMessageF(isQuitMessageF func(message tea.Msg) bool) Option {
	return func(r *Repl) {
		r.isQuitMessage = isQuitMessageF
	}
}
