package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
)

func RunInputLoop(paneID string) {
	historyFile := os.ExpandEnv("$HOME/.clp_history")

	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "You: ",
		HistoryFile: historyFile,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not start readline: %v\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			// Ctrl-D or EOF — exit cleanly, leave Claude pane running
			break
		}

		input := strings.TrimSpace(line)
		if input == "" {
			continue
		}

		switch input {
		case "/quit":
			// Exit input loop — Claude pane stays running full-width
			return
		case "/kill":
			// Kill the entire tmux session (both panes)
			sessionID, _, _ := strings.Cut(paneID, ":")
			exec.Command("tmux", "kill-session", "-t", sessionID).Run()
			return
		default:
			// Forward message to Claude pane as keystrokes
			exec.Command("tmux", "send-keys", "-t", paneID, input, "Enter").Run()
		}
	}
}
