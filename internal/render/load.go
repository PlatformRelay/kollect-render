package render

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Group is one stable groupBy bucket (key casing preserved from first seen item).
type Group struct {
	Key   string
	Items any
}

// LoadContextFile loads a RenderContext from a YAML or JSON file.
func LoadContextFile(path string) (RenderContext, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return RenderContext{}, fmt.Errorf("read %s: %w", path, err)
	}
	return DecodeContext(raw)
}

// DecodeContext decodes a YAML/JSON RenderContext document.
func DecodeContext(raw []byte) (RenderContext, error) {
	var ctx RenderContext
	if err := yaml.Unmarshal(raw, &ctx); err != nil {
		return RenderContext{}, fmt.Errorf("decode context: %w", err)
	}
	if ctx.Upstream == nil {
		ctx.Upstream = map[string]UpstreamEntry{}
	}
	if ctx.Copy == nil {
		ctx.Copy = map[string]string{}
	}
	return ctx, nil
}
