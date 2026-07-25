package render

import (
	"bytes"
	"fmt"
	"text/template"
)

// Render executes tmpl against ctx and returns deterministic bytes.
// Same RenderContext in ⇒ same bytes out (no clock/env reads).
func Render(tmpl string, ctx RenderContext) ([]byte, error) {
	t, err := template.New("render").Funcs(FuncMap()).Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}
	return buf.Bytes(), nil
}
