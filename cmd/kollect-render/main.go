package main

import (
	"fmt"
	"os"

	"github.com/platformrelay/kollect-render/internal/validate"
)

// version is a scaffold placeholder until release tagging lands (E2-S07).
const version = "0.0.0-dev"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		printUsage()
		return 2
	}
	switch args[0] {
	case "version":
		fmt.Println(version)
		return 0
	case "validate":
		return runValidate(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		return 2
	}
}

func runValidate(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: kollect-render validate <document>")
		return 2
	}
	if err := validate.File(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "validate: %v\n", err)
		return 2
	}
	return 0
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage: kollect-render <command>")
	fmt.Fprintln(os.Stderr, "commands: validate, version")
}
