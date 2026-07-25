// Package artifact defines the private-publisher render artifact contract (REQ-E2-S05-01):
// deterministic body bytes plus a digest/metadata sidecar. No network or credential access.
package artifact

import (
	"fmt"
)

// SchemaVersion pins the sidecar JSON shape for E5 consumers.
const SchemaVersion = "artifact-v0"

// SidecarSuffix is appended to the body path to form the companion metadata file.
const SidecarSuffix = ".meta.json"

// Meta is caller-supplied generation metadata recorded alongside the body digest.
type Meta struct {
	Format          string
	GeneratedAt     string // RFC3339 UTC
	Origin          string
	SnapshotSHA     string
	SourceRepoURL   string
	TemplateDigest  string // sha256:… of template bytes when --template is used
	RendererVersion string
}

// Sidecar is the digest/metadata companion consumed by the private E5 publisher.
type Sidecar struct {
	SchemaVersion   string `json:"schemaVersion"`
	ContentDigest   string `json:"contentDigest"`
	Format          string `json:"format"`
	ByteLength      int    `json:"byteLength"`
	GeneratedAt     string `json:"generatedAt,omitempty"`
	Origin          string `json:"origin,omitempty"`
	SnapshotSHA     string `json:"snapshotSHA,omitempty"`
	SourceRepoURL   string `json:"sourceRepoURL,omitempty"`
	TemplateDigest  string `json:"templateDigest,omitempty"`
	RendererVersion string `json:"rendererVersion,omitempty"`
}

// ContentDigest returns sha256:<hex> over body (stub until implemented).
func ContentDigest(body []byte) string {
	_ = body
	return ""
}

// Build constructs a sidecar for body + meta (stub until implemented).
func Build(body []byte, meta Meta) Sidecar {
	_ = body
	_ = meta
	return Sidecar{}
}

// MarshalSidecar returns deterministic JSON bytes for s (stub until implemented).
func MarshalSidecar(s Sidecar) ([]byte, error) {
	_ = s
	return nil, fmt.Errorf("artifact.MarshalSidecar: not implemented")
}

// SidecarPath returns the companion path for a body file (stub until implemented).
func SidecarPath(bodyPath string) string {
	_ = bodyPath
	return ""
}

// Write writes body to bodyPath and the sidecar to SidecarPath(bodyPath) (stub until implemented).
func Write(bodyPath string, body []byte, meta Meta) (Sidecar, error) {
	_ = bodyPath
	_ = body
	_ = meta
	return Sidecar{}, fmt.Errorf("artifact.Write: not implemented")
}
