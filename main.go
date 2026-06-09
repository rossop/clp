package main

import (
	"fmt"
	"os"
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
	fmt.Printf("startOrAttach: session=%s dir%s claudeArgs=%v\n", name, dir, claudeArgs)
}

func listSession() {
	fmt.Println("listSession: not yet implemented")
}

func killSession(id string) {
	fmt.Printf("killSession: %s\n", id)
}
