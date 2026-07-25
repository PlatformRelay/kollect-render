package render

import "fmt"

// Group is one stable groupBy bucket (key casing preserved from first seen item).
type Group struct {
	Key   string
	Items any
}

// LoadContextFile loads a RenderContext from a YAML or JSON file. Stub.
func LoadContextFile(_ string) (RenderContext, error) {
	return RenderContext{}, fmt.Errorf("LoadContextFile: not implemented")
}
