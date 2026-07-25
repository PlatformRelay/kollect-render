package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/platformrelay/kollect-render/internal/format"
	"github.com/platformrelay/kollect-render/internal/render"
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
	case "render":
		return runRender(args[1:])
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

func runRender(args []string) int {
	formatName := "markdown"
	templatePath := ""
	contextPath := ""
	outputPath := ""
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--format" && i+1 < len(args):
			i++
			formatName = args[i]
		case strings.HasPrefix(a, "--format="):
			formatName = strings.TrimPrefix(a, "--format=")
		case a == "--template" && i+1 < len(args):
			i++
			templatePath = args[i]
		case strings.HasPrefix(a, "--template="):
			templatePath = strings.TrimPrefix(a, "--template=")
		case a == "--context" && i+1 < len(args):
			i++
			contextPath = args[i]
		case strings.HasPrefix(a, "--context="):
			contextPath = strings.TrimPrefix(a, "--context=")
		case a == "--output" && i+1 < len(args):
			i++
			outputPath = args[i]
		case strings.HasPrefix(a, "--output="):
			outputPath = strings.TrimPrefix(a, "--output=")
		default:
			fmt.Fprintf(os.Stderr, "render: unknown argument %s\n", a)
			return 2
		}
	}
	if contextPath == "" {
		fmt.Fprintln(os.Stderr, "usage: kollect-render render --format <markdown|confluence-storage> --context <file> [--template <file>] [--output <file>]")
		return 2
	}
	enc, ok := format.Lookup(formatName)
	if !ok {
		fmt.Fprintf(os.Stderr, "render: unsupported format %q (registered: %s)\n", formatName, strings.Join(format.Names(), ", "))
		return 2
	}
	ctx, err := render.LoadContextFile(contextPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render: %v\n", err)
		return 2
	}
	var out []byte
	if templatePath != "" {
		tmplBytes, err := os.ReadFile(templatePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "render: read template: %v\n", err)
			return 2
		}
		out, err = render.Render(string(tmplBytes), ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "render: %v\n", err)
			return 2
		}
	} else {
		// Built-in env-inventory Model → encoder (REQ-E2-S04-01).
		model, err := format.EnvInventoryModel(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "render: %v\n", err)
			return 2
		}
		out, err = enc.Encode(model)
		if err != nil {
			fmt.Fprintf(os.Stderr, "render: encode %s: %v\n", enc.Name(), err)
			return 2
		}
	}
	if outputPath == "" || outputPath == "-" {
		if _, err := os.Stdout.Write(out); err != nil {
			fmt.Fprintf(os.Stderr, "render: write stdout: %v\n", err)
			return 2
		}
		return 0
	}
	if err := os.WriteFile(outputPath, out, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "render: write output: %v\n", err)
		return 2
	}
	return 0
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage: kollect-render <command>")
	fmt.Fprintln(os.Stderr, "commands: validate, render, version")
}
