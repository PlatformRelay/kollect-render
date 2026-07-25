package render

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

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

// FuncMap returns the template helpers bound for text/template execution.
func FuncMap() map[string]any {
	return map[string]any{
		"semverCompare":     SemverCompare,
		"versionDifference": VersionDifference,
		"statusEmoji":       StatusEmoji,
		"fmtTime":           FmtTime,
		"fmtDate":           FmtDate,
		"anchor":            Anchor,
		"sortBy":            SortBy,
		"groupBy":           GroupBy,
		"capNote":           CapNote,
	}
}

// SemverCompare returns -1/0/1 for a vs b; tolerant of v prefixes and build suffixes.
func SemverCompare(a, b string) int {
	ca, oa := canonicalize(a), canonicalize(b)
	if !semver.IsValid(ca) || !semver.IsValid(oa) {
		return strings.Compare(a, b)
	}
	return semver.Compare(ca, oa)
}

// VersionDifference classifies drift: same|minor|major|incomparable|unknown.
func VersionDifference(current, observed string) string {
	if strings.TrimSpace(current) == "" || strings.TrimSpace(observed) == "" {
		return "unknown"
	}
	c, o := canonicalize(current), canonicalize(observed)
	if !semver.IsValid(c) || !semver.IsValid(o) {
		return "incomparable"
	}
	if semver.Compare(c, o) == 0 {
		return "same"
	}
	if semver.Major(c) != semver.Major(o) {
		return "major"
	}
	return "minor"
}

// StatusEmoji maps a state token to its legend glyph (case-insensitive).
func StatusEmoji(state string) string {
	switch strings.ToLower(strings.TrimSpace(state)) {
	case "ready", "current":
		return "✅"
	case "warning", "minor", "stale":
		return "⚠️"
	case "notready", "major":
		return "🔴"
	default:
		return "◽"
	}
}

// FmtTime formats t as RFC3339 UTC.
func FmtTime(t any) string {
	ts, ok := asTime(t)
	if !ok {
		return ""
	}
	return ts.UTC().Format(time.RFC3339)
}

// FmtDate formats t as YYYY-MM-DD UTC.
func FmtDate(t any) string {
	ts, ok := asTime(t)
	if !ok {
		return ""
	}
	return ts.UTC().Format("2006-01-02")
}

var nonAnchor = regexp.MustCompile(`[^a-z0-9_]+`)

// Anchor returns a Confluence/markdown-safe anchor id.
func Anchor(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonAnchor.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return s
}

// SortBy stably sorts a slice of structs by field (case-insensitive string compare).
func SortBy(field string, items any) (any, error) {
	rv := reflect.ValueOf(items)
	if rv.Kind() != reflect.Slice {
		return nil, fmt.Errorf("sortBy: items must be a slice, got %T", items)
	}
	n := rv.Len()
	type pair struct {
		idx int
		key string
	}
	pairs := make([]pair, n)
	for i := 0; i < n; i++ {
		key, err := fieldString(rv.Index(i), field)
		if err != nil {
			return nil, err
		}
		pairs[i] = pair{idx: i, key: key}
	}
	sort.SliceStable(pairs, func(i, j int) bool {
		return strings.ToLower(pairs[i].key) < strings.ToLower(pairs[j].key)
	})
	out := reflect.MakeSlice(rv.Type(), n, n)
	for i, p := range pairs {
		out.Index(i).Set(rv.Index(p.idx))
	}
	return out.Interface(), nil
}

// GroupBy stably groups a slice of structs by field (case-insensitive).
// Group keys keep first-seen casing; groups are ordered by case-insensitive key.
func GroupBy(field string, items any) (any, error) {
	rv := reflect.ValueOf(items)
	if rv.Kind() != reflect.Slice {
		return nil, fmt.Errorf("groupBy: items must be a slice, got %T", items)
	}
	type bucket struct {
		key   string
		fold  string
		items []reflect.Value
	}
	order := make([]string, 0)
	byFold := map[string]*bucket{}
	for i := 0; i < rv.Len(); i++ {
		elem := rv.Index(i)
		key, err := fieldString(elem, field)
		if err != nil {
			return nil, err
		}
		fold := strings.ToLower(key)
		b, ok := byFold[fold]
		if !ok {
			b = &bucket{key: key, fold: fold}
			byFold[fold] = b
			order = append(order, fold)
		}
		b.items = append(b.items, elem)
	}
	sort.SliceStable(order, func(i, j int) bool {
		return order[i] < order[j]
	})
	groups := make([]Group, 0, len(order))
	elemType := rv.Type().Elem()
	for _, fold := range order {
		b := byFold[fold]
		slice := reflect.MakeSlice(reflect.SliceOf(elemType), len(b.items), len(b.items))
		for i, v := range b.items {
			slice.Index(i).Set(v)
		}
		groups = append(groups, Group{Key: b.key, Items: slice.Interface()})
	}
	return groups, nil
}

// CapNote renders the capped-list note when capped is true.
func CapNote(capped bool, n int) string {
	if !capped {
		return ""
	}
	return fmt.Sprintf("(list capped at %d entries)", n)
}

func canonicalize(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return v
	}
	if idx := strings.IndexByte(v, '+'); idx >= 0 {
		v = v[:idx]
	}
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	return v
}

func asTime(t any) (time.Time, bool) {
	switch v := t.(type) {
	case time.Time:
		return v, true
	case *time.Time:
		if v == nil {
			return time.Time{}, false
		}
		return *v, true
	case string:
		parsed, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return time.Time{}, false
		}
		return parsed, true
	default:
		return time.Time{}, false
	}
}

func fieldString(elem reflect.Value, field string) (string, error) {
	for elem.Kind() == reflect.Pointer {
		if elem.IsNil() {
			return "", nil
		}
		elem = elem.Elem()
	}
	if elem.Kind() != reflect.Struct {
		return "", fmt.Errorf("sortBy/groupBy: element must be struct, got %s", elem.Kind())
	}
	f := elem.FieldByName(field)
	if !f.IsValid() {
		return "", fmt.Errorf("sortBy/groupBy: no field %q on %s", field, elem.Type())
	}
	if f.Kind() == reflect.String {
		return f.String(), nil
	}
	return fmt.Sprint(f.Interface()), nil
}
