package format

import (
	"fmt"

	"github.com/platformrelay/kollect-render/internal/render"
)

// EnvInventoryModel builds the format-agnostic env-inventory page Model from a
// RenderContext. Same Model feeds every registered encoder (REQ-E2-S04-01).
// Stub until the builder lands.
func EnvInventoryModel(ctx render.RenderContext) (Model, error) {
	_ = ctx
	return Model{}, fmt.Errorf("format.EnvInventoryModel: not implemented")
}
