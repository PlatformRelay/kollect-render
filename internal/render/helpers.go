package render

import "fmt"

// HelperNames is the frozen helper set for the current minor (REQ-E2-S03-02).
// Adding a name is a minor change; removing/re-typing is breaking.
var HelperNames = []string{
	"semverCompare",
	"versionDifference",
	"statusEmoji",
	"fmtTime",
	"fmtDate",
	"anchor",
	"sortBy",
	"groupBy",
	"capNote",
}

// SemverCompare returns -1/0/1 for a vs b. Stub until engine lands.
func SemverCompare(a, b string) int {
	_ = a
	_ = b
	return 99
}

// VersionDifference classifies drift between current and observed. Stub.
func VersionDifference(current, observed string) string {
	_ = current
	_ = observed
	return "stub"
}

// StatusEmoji maps a state token to its legend glyph. Stub.
func StatusEmoji(state string) string {
	_ = state
	return "?"
}

// FmtTime formats t as RFC3339 UTC. Stub.
func FmtTime(t any) string {
	_ = t
	return "stub-time"
}

// FmtDate formats t as YYYY-MM-DD UTC. Stub.
func FmtDate(t any) string {
	_ = t
	return "stub-date"
}

// Anchor returns a Confluence/markdown-safe anchor id. Stub.
func Anchor(s string) string {
	_ = s
	return "stub-anchor"
}

// SortBy stably sorts items by field (case-insensitive). Stub.
func SortBy(field string, items any) (any, error) {
	_ = field
	return items, fmt.Errorf("sortBy: not implemented")
}

// GroupBy stably groups items by field (case-insensitive). Stub.
func GroupBy(field string, items any) (any, error) {
	_ = field
	return nil, fmt.Errorf("groupBy: not implemented")
}

// CapNote renders the capped-list note when capped is true. Stub.
func CapNote(capped bool, n int) string {
	_ = capped
	_ = n
	return "stub-cap"
}
