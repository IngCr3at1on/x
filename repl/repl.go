package repl

import (
	"bytes"
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	Repl struct {
		ctx  context.Context
		m    model
		buf  *bytes.Buffer // FIXME: replace with append only read/writer that removes older lines after reaching max.
		prog *tea.Program
		eval EvalFunc

		prompt        string
		isQuitMessage func(message tea.Msg) bool

		ready    bool
		highPerf bool

		errCh chan error
	}

	EvalFunc func(ctx context.Context, txt string, r *Repl) (tea.Cmd, error)
)

// Start a Repl with the provided EvalFunc and Options.
// EvalFunc is called on enter.
func Start(eval EvalFunc, opts ...Option) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := &Repl{
		ctx:           ctx,
		buf:           new(bytes.Buffer),
		eval:          eval,
		prompt:        defaultPrompt,
		isQuitMessage: defaultIsQuitMessageF,
		// FIXME: this is forced to on because for some reason it without I have vanishing text in my viewport.
		highPerf: true,
		errCh:    make(chan error, 1),
	}

	for _, opt := range opts {
		opt(r)
	}

	r.m = initialModel(r) // FIXME: set with option.
	r.prog = tea.NewProgram(r.m,
		tea.WithAltScreen())

	go func() {
		if err := r.prog.Start(); err != nil {
			r.errCh <- err
		}
		close(r.errCh)
	}()

	err := <-r.errCh
	if err != nil {
		return err
	}

	return nil
}

// ToViewport writes to the internal buffer used to hold content for the viewport.
// This is expected to be called from inside of an EvalFunc, calling this from
// outside of an EvalFunc will cause it not to render until a tea.WindowSizeMsg
// or tea.KeyEnter is received within model.Update.
func (r Repl) ToViewport(txt string) error {
	_, err := r.buf.WriteString(txt + "\n")
	if err != nil {
		return err
	}

	return nil
}

func (r Repl) error(err error) {
	r.prog.Kill()
	r.errCh <- err
	close(r.errCh)
}
