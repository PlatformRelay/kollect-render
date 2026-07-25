package format

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// ValidateStorageFragment checks that a confluence-storage body is basically
// well-formed XHTML (REQ-E2-S04-02). Confluence prefixes (ac:) are declared on
// a wrapper so standard XML tokenization applies.
func ValidateStorageFragment(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("storage fragment: empty")
	}
	wrapped := `<storage xmlns:ac="https://atlassian.com/content">` + s + `</storage>`
	dec := xml.NewDecoder(strings.NewReader(wrapped))
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("storage fragment: %w", err)
		}
		_ = tok
	}
}
