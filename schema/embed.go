// Package schema embeds the committed inventory JSON Schema documents.
//
// InventoryV0 is the UTF-8 bytes of inventory-v0.schema.json - the draft v0
// neutral inventory evidence envelope schema (JSON Schema draft 2020-12).
// Documents validated against it must set schemaVersion to the const pinned in
// that file (currently "0.2.0"). This is the module's only intentionally
// importable surface; render/validate/format/artifact live under internal/.
package schema

// Blank import enables //go:embed below; InventoryV0 is []byte, not embed.FS,
// so the package is never referenced by name.
import _ "embed"

// InventoryV0 is the UTF-8 contents of inventory-v0.schema.json, the draft v0
// inventory evidence envelope schema this module pins for validation.
//
//go:embed inventory-v0.schema.json
var InventoryV0 []byte
