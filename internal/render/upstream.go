package render

import (
	"fmt"
	"time"
)

// LoadUpstreamFile loads a map[componentID]UpstreamEntry YAML file (REQ-E2-S06-01).
// Stub — implemented in the S06 feat commit.
func LoadUpstreamFile(path string) (map[string]UpstreamEntry, error) {
	_ = path
	return nil, fmt.Errorf("LoadUpstreamFile: not implemented")
}

// ApplyGenerationOverrides applies optional --generated-at / --report-origin overrides.
// Stub — implemented in the S06 feat commit.
func ApplyGenerationOverrides(ctx RenderContext, generatedAt, reportOrigin string) (RenderContext, error) {
	_ = time.Time{}
	_ = generatedAt
	_ = reportOrigin
	return ctx, fmt.Errorf("ApplyGenerationOverrides: not implemented")
}
