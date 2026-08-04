// Package artifact defines the private-publisher render artifact contract
// (REQ-E2-S05-01): deterministic body bytes plus a digest/metadata sidecar
// (*.meta.json). The package never opens the network or holds credentials.
package artifact

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
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

// ContentDigest returns sha256:<hex> over body.
func ContentDigest(body []byte) string {
	if body == nil {
		body = []byte{}
	}
	sum := sha256.Sum256(body)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// Build constructs a sidecar for body + meta.
func Build(body []byte, meta Meta) Sidecar {
	if body == nil {
		body = []byte{}
	}
	return Sidecar{
		SchemaVersion:   SchemaVersion,
		ContentDigest:   ContentDigest(body),
		Format:          meta.Format,
		ByteLength:      len(body),
		GeneratedAt:     meta.GeneratedAt,
		Origin:          meta.Origin,
		SnapshotSHA:     meta.SnapshotSHA,
		SourceRepoURL:   meta.SourceRepoURL,
		TemplateDigest:  meta.TemplateDigest,
		RendererVersion: meta.RendererVersion,
	}
}

// MarshalSidecar returns deterministic JSON bytes for s (2-space indent + trailing newline).
func MarshalSidecar(s Sidecar) ([]byte, error) {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal sidecar: %w", err)
	}
	return append(b, '\n'), nil
}

// SidecarPath returns the companion metadata path for a body file.
func SidecarPath(bodyPath string) string {
	return bodyPath + SidecarSuffix
}

// Write writes body to bodyPath and the sidecar to SidecarPath(bodyPath).
func Write(bodyPath string, body []byte, meta Meta) (Sidecar, error) {
	if bodyPath == "" {
		return Sidecar{}, fmt.Errorf("artifact.Write: empty body path")
	}
	sc := Build(body, meta)
	raw, err := MarshalSidecar(sc)
	if err != nil {
		return Sidecar{}, err
	}
	if err := os.WriteFile(bodyPath, body, 0o644); err != nil {
		return Sidecar{}, fmt.Errorf("write body: %w", err)
	}
	if err := os.WriteFile(SidecarPath(bodyPath), raw, 0o644); err != nil {
		return Sidecar{}, fmt.Errorf("write sidecar: %w", err)
	}
	return sc, nil
}
