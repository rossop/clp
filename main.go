package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const sessionPrefix = "clp-"

func main() {
	args := os.Args[1:]

	// Internal subcommand: called by clp itself inside the input pane
	if len(args) >= 2 && args[0] == "--input-loop" {
		RunInputLoop(args[1])
		return
	}

	if len(args) == 0 {
		startOrAttach("default", "", nil)
		return
	}

	switch args[0] {
	case "list":
		listSession()
	case "kill":
		name := "default"
		if len(args) > 1 {
			name = args[1]
		}
		killSession(sessionPrefix + name)
	default:
		startOrAttach(args[0], "", nil)
	}
}

func startOrAttach(name, dir string, claudeArgs []string) {
	id := sessionPrefix + name
	if sessionExists(id) {
		attach(id)
		return
	}
	createSession(id, dir, claudeArgs)
}

func sessionExists(id string) bool {
	return exec.Command("tmux", "has-session", "-t", id).Run() == nil
}

func createSession(id, dir string, claudeArgs []string) {
	// 1. Create detached tmux session, optionally in a working directory
	newArgs := []string{"new-session", "-d", "-s", id}
	if dir != "" {
		newArgs = append(newArgs, "-c", dir)
	}
	if err := exec.Command("tmux", newArgs...).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: could not create tmux session: %v\n", err)
		os.Exit(1)
	}

	// 2. Start claude in the left pane (pane 0)
	claudeCmd := strings.Join(append([]string{"claude"}, claudeArgs...), " ")
	exec.Command("tmux", "send-keys", "-t", id+":0.0", claudeCmd, "Enter").Run()

	// 3. Split right pane at 30% width
	exec.Command("tmux", "split-window", "-h", "-p", "30", "-t", id+":0.0").Run()

	// 4. Start input loop in the right pane, pointing it at the claude pane
	exe, _ := os.Executable()
	exec.Command("tmux", "send-keys", "-t", id+":0.1",
		exe+" --input-loop "+id+":0.0", "Enter").Run()

	// 5. Focus on the input pane so the user lands there on attach
	exec.Command("tmux", "select-pane", "-t", id+":0.1").Run()

	// 6. Replace this process with tmux attach. No orphan parent
	attach(id)
}

func attach(id string) {
	tmux, err := exec.LookPath("tmux")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: tmux not found in PATH")
		os.Exit(1)
	}
	syscall.Exec(tmux, []string{"tmux", "attach-session", "-t", id}, os.Environ())
}

func listSession() {
	out, err := exec.Command("tmux", "list-sessions", "-F", "#{session_name}").Output()
	if err != nil {
		fmt.Println("no active clp sessions")
		return
	}
	found := false
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		if name, ok := strings.CutPrefix(line, sessionPrefix); ok {
			fmt.Println(name)
			found = true
		}
	}
	if !found {
		fmt.Println("no active clp sessions")
	}
}

func killSession(id string) {
	if err := exec.Command("tmux", "kill-session", "-t", id).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "no session named: %s\n", strings.TrimPrefix(id, sessionPrefix))
		os.Exit(1)
	}
	fmt.Printf("killed: %s\n", strings.TrimPrefix(id, sessionPrefix))
}
