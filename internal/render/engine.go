package render

import "fmt"

// Render executes tmpl against ctx and returns deterministic bytes.
// Stub: returns a hard failure until the template engine lands (REQ-E2-S03-01).
func Render(tmpl string, ctx RenderContext) ([]byte, error) {
	_ = tmpl
	_ = ctx
	return nil, fmt.Errorf("render: not implemented")
}
