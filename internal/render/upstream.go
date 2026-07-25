package render

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// LoadUpstreamFile loads a map[componentID]UpstreamEntry YAML document
// (the shape under RenderContext.Upstream / template-contract §1). Used by
// --upstream-deps (REQ-E2-S06-01). Completeness/catalog policy stays with the
// private aggregator — this only decodes the already-shaped map.
func LoadUpstreamFile(path string) (map[string]UpstreamEntry, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read upstream-deps %s: %w", path, err)
	}
	var m map[string]UpstreamEntry
	if err := yaml.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("decode upstream-deps %s: %w", path, err)
	}
	if m == nil {
		m = map[string]UpstreamEntry{}
	}
	return m, nil
}

// ApplyGenerationOverrides applies optional --generated-at / --report-origin
// overrides onto ctx.Generation. Empty strings leave the corresponding field
// unchanged. generatedAt must be RFC3339 (UTC) when non-empty.
func ApplyGenerationOverrides(ctx RenderContext, generatedAt, reportOrigin string) (RenderContext, error) {
	if generatedAt != "" {
		t, err := time.Parse(time.RFC3339, generatedAt)
		if err != nil {
			return RenderContext{}, fmt.Errorf("invalid --generated-at %q: %w", generatedAt, err)
		}
		ctx.Generation.GeneratedAt = t.UTC()
	}
	if reportOrigin != "" {
		ctx.Generation.Origin = reportOrigin
	}
	return ctx, nil
}
