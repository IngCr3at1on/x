package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ingcr3at1on/x/repl"
)

func echo(_ context.Context, txt string, r *repl.Repl) (tea.Cmd, error) {
	switch txt {
	case "quit", "exit":
		return tea.Quit, nil
	default:
		return nil, r.ToViewport(txt)
	}
}

func main() {
	if err := repl.Start(echo); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
