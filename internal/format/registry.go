package format

// Well-known format names (REQ-E2-S04-01).
const (
	NameMarkdown          = "markdown"
	NameConfluenceStorage = "confluence-storage"
)

// Encoder turns a Model into format-specific bytes.
type Encoder interface {
	Name() string
	Encode(m Model) ([]byte, error)
}

var registry = map[string]Encoder{}

func init() {
	MustRegister(markdownEncoder{})
	MustRegister(confluenceEncoder{})
}

// MustRegister adds an encoder to the compile-time registry. Panics on duplicate.
func MustRegister(e Encoder) {
	name := e.Name()
	if _, ok := registry[name]; ok {
		panic("format: duplicate registration: " + name)
	}
	registry[name] = e
}

// Lookup returns the encoder for name, or false if unknown.
func Lookup(name string) (Encoder, bool) {
	e, ok := registry[name]
	return e, ok
}

// Names returns registered format names in stable order.
func Names() []string {
	out := make([]string, 0, len(registry))
	// Stable: markdown then confluence-storage then any others sorted.
	prefer := []string{NameMarkdown, NameConfluenceStorage}
	seen := map[string]bool{}
	for _, n := range prefer {
		if _, ok := registry[n]; ok {
			out = append(out, n)
			seen[n] = true
		}
	}
	rest := make([]string, 0)
	for n := range registry {
		if !seen[n] {
			rest = append(rest, n)
		}
	}
	for i := 0; i < len(rest); i++ {
		for j := i + 1; j < len(rest); j++ {
			if rest[j] < rest[i] {
				rest[i], rest[j] = rest[j], rest[i]
			}
		}
	}
	return append(out, rest...)
}
