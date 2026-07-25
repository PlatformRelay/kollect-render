package main

import (
	"fmt"
	"os"
)

// version is a scaffold placeholder until release tagging lands (E2-S07).
const version = "0.0.0-dev"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: kollect-render <command>")
		fmt.Fprintln(os.Stderr, "commands: version")
		return 2
	}
	switch args[0] {
	case "version":
		fmt.Println(version)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		return 2
	}
}
