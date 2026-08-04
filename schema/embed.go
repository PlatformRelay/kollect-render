// Package schema embeds the committed inventory JSON Schema documents.
package schema

// Blank import enables //go:embed below; InventoryV0 is []byte, not embed.FS,
// so the package is never referenced by name.
import _ "embed"

// InventoryV0 is the draft v0 inventory evidence envelope schema.
//
//go:embed inventory-v0.schema.json
var InventoryV0 []byte
