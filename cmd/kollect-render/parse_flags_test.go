package main

import (
	"strings"
	"testing"
)

func TestParseRenderFlagsSpaceSeparated(t *testing.T) {
	t.Parallel()
	opts, err := parseRenderFlags([]string{
		"--format", "confluence-storage",
		"--context", "/tmp/ctx.yaml",
		"--template", "/tmp/t.md",
		"--output", "/tmp/out.md",
		"--upstream-deps", "/tmp/deps.yaml",
		"--generated-at", "2026-07-25T12:00:00Z",
		"--report-origin", "ci",
	})
	if err != nil {
		t.Fatalf("parseRenderFlags: %v", err)
	}
	if opts.formatName != "confluence-storage" {
		t.Fatalf("formatName = %q", opts.formatName)
	}
	if opts.contextPath != "/tmp/ctx.yaml" || opts.templatePath != "/tmp/t.md" {
		t.Fatalf("paths = %+v", opts)
	}
	if opts.outputPath != "/tmp/out.md" || opts.upstreamDepsPath != "/tmp/deps.yaml" {
		t.Fatalf("output/deps = %+v", opts)
	}
	if opts.generatedAt != "2026-07-25T12:00:00Z" || opts.reportOrigin != "ci" {
		t.Fatalf("overrides = %+v", opts)
	}
}

func TestParseRenderFlagsEqualsForm(t *testing.T) {
	t.Parallel()
	opts, err := parseRenderFlags([]string{
		"--format=markdown",
		"--context=/tmp/ctx.yaml",
		"--output=/tmp/out.md",
		"--generated-at=2026-07-25T00:00:00Z",
		"--report-origin=manual",
	})
	if err != nil {
		t.Fatalf("parseRenderFlags: %v", err)
	}
	if opts.formatName != "markdown" || opts.contextPath != "/tmp/ctx.yaml" {
		t.Fatalf("opts = %+v", opts)
	}
	if opts.outputPath != "/tmp/out.md" {
		t.Fatalf("outputPath = %q", opts.outputPath)
	}
	if opts.generatedAt != "2026-07-25T00:00:00Z" || opts.reportOrigin != "manual" {
		t.Fatalf("overrides = %+v", opts)
	}
}

func TestParseRenderFlagsTrailingFormatNoValue(t *testing.T) {
	t.Parallel()
	_, err := parseRenderFlags([]string{"--context", "/tmp/ctx.yaml", "--format"})
	if err == nil {
		t.Fatal("expected error for trailing --format with no value")
	}
	if got := err.Error(); got != "render: unknown argument --format" {
		t.Fatalf("error = %q, want unknown argument --format", got)
	}
}

func TestParseRenderFlagsUnknownArgument(t *testing.T) {
	t.Parallel()
	_, err := parseRenderFlags([]string{"--nope"})
	if err == nil {
		t.Fatal("expected error for --nope")
	}
	if got := err.Error(); got != "render: unknown argument --nope" {
		t.Fatalf("error = %q", got)
	}
}

func TestParseRenderFlagsDuplicateLastWins(t *testing.T) {
	t.Parallel()
	opts, err := parseRenderFlags([]string{
		"--format", "markdown",
		"--format=confluence-storage",
		"--context", "/tmp/a.yaml",
		"--context=/tmp/b.yaml",
		"--output", "/tmp/old.md",
		"--output", "/tmp/new.md",
	})
	if err != nil {
		t.Fatalf("parseRenderFlags: %v", err)
	}
	if opts.formatName != "confluence-storage" {
		t.Fatalf("duplicate format last-wins: got %q", opts.formatName)
	}
	if opts.contextPath != "/tmp/b.yaml" {
		t.Fatalf("duplicate context last-wins: got %q", opts.contextPath)
	}
	if opts.outputPath != "/tmp/new.md" {
		t.Fatalf("duplicate output last-wins: got %q", opts.outputPath)
	}
}

func TestParseRenderFlagsDefaultFormat(t *testing.T) {
	t.Parallel()
	opts, err := parseRenderFlags([]string{"--context", "/tmp/ctx.yaml"})
	if err != nil {
		t.Fatalf("parseRenderFlags: %v", err)
	}
	if opts.formatName != "markdown" {
		t.Fatalf("default format = %q, want markdown", opts.formatName)
	}
	if strings.TrimSpace(opts.templatePath) != "" {
		t.Fatalf("templatePath should be empty, got %q", opts.templatePath)
	}
}
