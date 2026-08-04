// Package main is the kollect-render CLI: credential-free validate and render
// of inventory documents into deterministic publisher artifacts.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/platformrelay/kollect-render/internal/artifact"
	"github.com/platformrelay/kollect-render/internal/format"
	"github.com/platformrelay/kollect-render/internal/render"
	"github.com/platformrelay/kollect-render/internal/validate"
)

// version is overridden by goreleaser ldflags (-X main.version={{.Version}}).
var version = "0.0.0-dev"

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

// renderOptions holds parsed `render` CLI flags.
type renderOptions struct {
	formatName       string
	templatePath     string
	contextPath      string
	outputPath       string
	upstreamDepsPath string
	generatedAt      string
	reportOrigin     string
}

// renderFlagSetters maps long-option names to option fields. Kept as a package
// var so parseRenderFlags stays a thin loop under the gocyclo floor.
var renderFlagSetters = map[string]func(*renderOptions, string){
	"--format":        func(o *renderOptions, v string) { o.formatName = v },
	"--template":      func(o *renderOptions, v string) { o.templatePath = v },
	"--context":       func(o *renderOptions, v string) { o.contextPath = v },
	"--output":        func(o *renderOptions, v string) { o.outputPath = v },
	"--upstream-deps": func(o *renderOptions, v string) { o.upstreamDepsPath = v },
	"--generated-at":  func(o *renderOptions, v string) { o.generatedAt = v },
	"--report-origin": func(o *renderOptions, v string) { o.reportOrigin = v },
}

// parseRenderFlags parses render subcommand flags. Supports `--flag value` and
// `--flag=value`. Duplicate flags: last wins. A bare `--flag` with no value is
// treated as an unknown argument (same as the historical CLI).
func parseRenderFlags(args []string) (renderOptions, error) {
	opts := renderOptions{formatName: "markdown"}
	for i := 0; i < len(args); i++ {
		a := args[i]
		name, value, ok := splitRenderFlag(a)
		if ok {
			set, known := renderFlagSetters[name]
			if !known {
				return renderOptions{}, fmt.Errorf("render: unknown argument %s", a)
			}
			set(&opts, value)
			continue
		}
		set, known := renderFlagSetters[a]
		rest := args[i+1:]
		if !known || len(rest) == 0 {
			return renderOptions{}, fmt.Errorf("render: unknown argument %s", a)
		}
		i++
		set(&opts, rest[0])
	}
	return opts, nil
}

// splitRenderFlag parses `--name=value`. Returns ok=false when a is not equals-form.
func splitRenderFlag(a string) (name, value string, ok bool) {
	if !strings.HasPrefix(a, "--") {
		return "", "", false
	}
	eq := strings.IndexByte(a, '=')
	if eq <= 2 {
		return "", "", false
	}
	return a[:eq], a[eq+1:], true
}

func runRender(args []string) int {
	opts, err := parseRenderFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	if opts.contextPath == "" {
		fmt.Fprintln(os.Stderr, "usage: kollect-render render --format <markdown|confluence-storage> --context <file> [--template <file>] [--output <file>] [--upstream-deps <file>] [--generated-at <RFC3339>] [--report-origin <label>]")
		fmt.Fprintln(os.Stderr, "note: --template is markdown-only; omit it to encode via the format registry (Model → encoder)")
		return 2
	}
	enc, ok := format.Lookup(opts.formatName)
	if !ok {
		fmt.Fprintf(os.Stderr, "render: unsupported format %q (registered: %s)\n", opts.formatName, strings.Join(format.Names(), ", "))
		return 2
	}
	// Templates are text/template markdown sources. Non-markdown formats must use the
	// registry Encode path (no --template); never silently ignore --format.
	if opts.templatePath != "" && opts.formatName != format.NameMarkdown {
		fmt.Fprintf(os.Stderr, "render: --template requires --format %s (got %q); omit --template to encode %s via the registry\n",
			format.NameMarkdown, opts.formatName, opts.formatName)
		return 2
	}
	ctx, err := loadRenderContext(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render: %v\n", err)
		return 2
	}
	out, templateDigest, err := produceRenderBytes(opts, enc, ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	return emitRender(opts.outputPath, out, artifact.Meta{
		Format:          opts.formatName,
		GeneratedAt:     render.FmtTime(ctx.Generation.GeneratedAt),
		Origin:          ctx.Generation.Origin,
		SnapshotSHA:     ctx.Generation.SnapshotSHA,
		SourceRepoURL:   ctx.Generation.SourceRepoURL,
		TemplateDigest:  templateDigest,
		RendererVersion: version,
	})
}

func loadRenderContext(opts renderOptions) (render.RenderContext, error) {
	ctx, err := render.LoadContextFile(opts.contextPath)
	if err != nil {
		return render.RenderContext{}, err
	}
	if opts.upstreamDepsPath != "" {
		up, err := render.LoadUpstreamFile(opts.upstreamDepsPath)
		if err != nil {
			return render.RenderContext{}, err
		}
		ctx.Upstream = up
	}
	return render.ApplyGenerationOverrides(ctx, opts.generatedAt, opts.reportOrigin)
}

func produceRenderBytes(opts renderOptions, enc format.Encoder, ctx render.RenderContext) ([]byte, string, error) {
	if opts.templatePath != "" {
		tmplBytes, err := os.ReadFile(opts.templatePath)
		if err != nil {
			return nil, "", fmt.Errorf("render: read template: %w", err)
		}
		sum := sha256.Sum256(tmplBytes)
		digest := "sha256:" + hex.EncodeToString(sum[:])
		out, err := render.Render(string(tmplBytes), ctx)
		if err != nil {
			return nil, "", fmt.Errorf("render: %w", err)
		}
		return out, digest, nil
	}
	// Built-in env-inventory Model → encoder (REQ-E2-S04-01).
	model, err := format.EnvInventoryModel(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("render: %w", err)
	}
	out, err := enc.Encode(model)
	if err != nil {
		return nil, "", fmt.Errorf("render: encode %s: %w", enc.Name(), err)
	}
	return out, "", nil
}

// emitRender writes rendered bytes to stdout ("-" / empty) or an artifact file + sidecar.
func emitRender(outputPath string, out []byte, meta artifact.Meta) int {
	if outputPath == "" || outputPath == "-" {
		if _, err := os.Stdout.Write(out); err != nil {
			fmt.Fprintf(os.Stderr, "render: write stdout: %v\n", err)
			return 2
		}
		return 0
	}
	// File output: body + digest/metadata sidecar for the private publisher (REQ-E2-S05-01).
	if _, err := artifact.Write(outputPath, out, meta); err != nil {
		fmt.Fprintf(os.Stderr, "render: write artifact: %v\n", err)
		return 2
	}
	return 0
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage: kollect-render <command>")
	fmt.Fprintln(os.Stderr, "commands: validate, render, version")
}
