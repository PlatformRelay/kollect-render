package format

import (
	"fmt"
	"strings"

	"github.com/platformrelay/kollect-render/internal/render"
)

// EnvInventoryModel builds the format-agnostic env-inventory page Model from a
// RenderContext. The same Model feeds every registered encoder (REQ-E2-S04-01).
func EnvInventoryModel(ctx render.RenderContext) (Model, error) {
	howTo, _ := ctx.Copy["how-to-change"]
	blocks := []Block{
		Heading{Level: 1, Text: "Environment inventory"},
		Blank{},
		Banner{
			GeneratedAt:   render.FmtTime(ctx.Generation.GeneratedAt),
			GeneratedDate: render.FmtDate(ctx.Generation.GeneratedAt),
			Origin:        ctx.Generation.Origin,
			SourceRepoURL: ctx.Generation.SourceRepoURL,
			SnapshotSHA:   ctx.Generation.SnapshotSHA,
			HowToChange:   howTo,
		},
		Blank{},
		StatusLegend{},
		Blank{},
		Heading{Level: 2, Text: "Sources"},
	}

	srcItems := make([][]Inline, 0, len(ctx.Manifest.Sources))
	for _, s := range ctx.Manifest.Sources {
		srcItems = append(srcItems, []Inline{
			Code{S: s.SourceID},
			Text{S: " — " + s.Status + " "},
			Emoji{S: render.StatusEmoji(s.Status)},
		})
	}
	blocks = append(blocks, BulletList{Items: srcItems}, Blank{}, Heading{Level: 2, Text: "Components"})

	for _, doc := range ctx.Docs {
		anchorID := render.Anchor(fmt.Sprintf("%s-%s", doc.Metadata.EnvironmentID, doc.Metadata.SourceID))
		blocks = append(blocks, Heading{
			Level:  3,
			Text:   doc.Metadata.EnvironmentID + " / " + doc.Metadata.SourceID,
			Anchor: anchorID,
		})
		blocks = append(blocks, Paragraph{Inlines: []Inline{
			Text{S: "Collected: " + render.FmtTime(doc.Metadata.CollectedAt)},
		}})

		comps, err := render.SortBy("ComponentID", doc.Components)
		if err != nil {
			return Model{}, err
		}
		compSlice, ok := comps.([]render.Component)
		if !ok {
			return Model{}, fmt.Errorf("format: unexpected components type %T", comps)
		}
		compItems := make([][]Inline, 0, len(compSlice))
		for _, c := range compSlice {
			item := []Inline{
				Strong{S: c.Name},
				Text{S: " ("},
				Code{S: c.ComponentID},
				Text{S: "): " + c.Version},
			}
			if up, ok := ctx.Upstream[c.ComponentID]; ok {
				item = append(item,
					Text{S: " · upstream " + up.ObservedVersion + " (" + render.VersionDifference(c.Version, up.ObservedVersion) + ") · status " + up.Status},
				)
			}
			compItems = append(compItems, item)
		}
		if len(compItems) > 0 {
			blocks = append(blocks, BulletList{Items: compItems})
		}

		helms, err := render.SortBy("Name", doc.HelmReleases)
		if err != nil {
			return Model{}, err
		}
		helmSlice, ok := helms.([]render.HelmRelease)
		if !ok {
			return Model{}, fmt.Errorf("format: unexpected helm type %T", helms)
		}
		helmItems := make([][]Inline, 0, len(helmSlice))
		for _, h := range helmSlice {
			helmItems = append(helmItems, []Inline{
				Text{S: "helm "},
				Code{S: h.Name},
				Text{S: " " + h.ChartVersion + " (" + h.Status + ")"},
			})
		}
		if len(helmItems) > 0 {
			blocks = append(blocks, Blank{}, BulletList{Items: helmItems}, Blank{})
		} else {
			blocks = append(blocks, Blank{}, Blank{})
		}

		groups, err := render.GroupBy("Namespace", doc.Workloads)
		if err != nil {
			return Model{}, err
		}
		groupSlice, ok := groups.([]render.Group)
		if !ok {
			return Model{}, fmt.Errorf("format: unexpected group type %T", groups)
		}
		for _, g := range groupSlice {
			blocks = append(blocks, Heading{Level: 4, Text: "Workloads in `" + g.Key + "`"})
			wlSorted, err := render.SortBy("Name", g.Items)
			if err != nil {
				return Model{}, err
			}
			wlSlice, ok := wlSorted.([]render.Workload)
			if !ok {
				return Model{}, fmt.Errorf("format: unexpected workload type %T", wlSorted)
			}
			wlItems := make([][]Inline, 0, len(wlSlice))
			for _, w := range wlSlice {
				wlItems = append(wlItems, []Inline{
					Text{S: w.Kind + "/" + w.Name + ": " + strings.Join(w.Images, ", ")},
				})
			}
			blocks = append(blocks, BulletList{Items: wlItems}, Blank{})
		}

		kube := strings.Join(doc.Nodes.KubeVersions, ", ")
		blocks = append(blocks, Paragraph{Inlines: []Inline{
			Text{S: fmt.Sprintf("Nodes: %d (%s)", doc.Nodes.Count, kube)},
		}})

		tgtItems := make([][]Inline, 0, len(doc.TargetOutcomes))
		for _, t := range doc.TargetOutcomes {
			line := []Inline{
				Text{S: "target "},
				Code{S: t.TargetID},
				Text{S: ": " + t.Status},
			}
			if t.ErrorCode != "" {
				line = append(line, Text{S: " (" + t.ErrorCode + ")"})
			}
			tgtItems = append(tgtItems, line)
		}
		if len(tgtItems) > 0 {
			blocks = append(blocks, BulletList{Items: tgtItems})
		}
		blocks = append(blocks, Blank{})
	}

	blocks = append(blocks,
		CapNote{N: 3},
		Blank{},
		Footnote{Text: "Health values reflect state at collection time — documentation, not monitoring."},
	)
	return Model{Blocks: blocks}, nil
}
