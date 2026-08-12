// Package validate checks inventory documents against schema.InventoryV0
// (draft v0 evidence envelope). YAML or JSON input is accepted; failures are
// reported as JSON Pointer-located schema errors with no network access.
package validate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/platformrelay/kollect-render/schema"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"
)

const schemaURL = "https://github.com/platformrelay/kollect-render/schema/inventory-v0.schema.json"

var schemaPrinter = message.NewPrinter(language.English)

// Error is a single schema violation located by JSON Pointer.
type Error struct {
	Pointer string
	Message string
}

func (e Error) Error() string {
	if e.Pointer == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Pointer, e.Message)
}

// Result holds zero or more pointered validation errors.
type Result struct {
	Errors []Error
}

func (r *Result) Error() string {
	if r == nil || len(r.Errors) == 0 {
		return ""
	}
	if len(r.Errors) == 1 {
		return r.Errors[0].Error()
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d validation errors:", len(r.Errors))
	for _, e := range r.Errors {
		b.WriteByte('\n')
		b.WriteString(e.Error())
	}
	return b.String()
}

var (
	schemaOnce sync.Once
	schemaObj  *jsonschema.Schema
	schemaErr  error
)

// Bytes validates a YAML or JSON inventory document against schema v0.
func Bytes(raw []byte) error {
	sch, err := loadSchema()
	if err != nil {
		return err
	}
	value, err := decodeDocument(raw)
	if err != nil {
		return err
	}
	if err := sch.Validate(value); err != nil {
		return pointeredResult(err)
	}
	return nil
}

// File reads path and validates it as an inventory document.
func File(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	return Bytes(raw)
}

func decodeDocument(raw []byte) (any, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, &Result{Errors: []Error{{Pointer: "", Message: "document is empty"}}}
	}
	var asYAML any
	if err := yaml.Unmarshal(raw, &asYAML); err != nil {
		return nil, fmt.Errorf("parse document: %w", err)
	}
	// jsonschema needs JSON number/bool types; re-encode via encoding/json.
	jb, err := json.Marshal(asYAML)
	if err != nil {
		return nil, fmt.Errorf("json encode: %w", err)
	}
	var v any
	if err := json.Unmarshal(jb, &v); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}
	return v, nil
}

func loadSchema() (*jsonschema.Schema, error) {
	schemaOnce.Do(func() {
		c := jsonschema.NewCompiler()
		c.AssertFormat()
		doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(schema.InventoryV0))
		if err != nil {
			schemaErr = fmt.Errorf("decode schema: %w", err)
			return
		}
		if err := c.AddResource(schemaURL, doc); err != nil {
			schemaErr = fmt.Errorf("add schema resource: %w", err)
			return
		}
		schemaObj, schemaErr = c.Compile(schemaURL)
		if schemaErr != nil {
			schemaErr = fmt.Errorf("compile schema: %w", schemaErr)
		}
	})
	return schemaObj, schemaErr
}

func pointeredResult(err error) error {
	ve, ok := err.(*jsonschema.ValidationError)
	if !ok {
		return err
	}
	res := &Result{}
	collectLeaves(ve, res)
	if len(res.Errors) == 0 {
		res.Errors = append(res.Errors, Error{Pointer: pointerFor(ve), Message: validationMessage(ve)})
	}
	return res
}

func collectLeaves(ve *jsonschema.ValidationError, res *Result) {
	if ve == nil {
		return
	}
	if len(ve.Causes) == 0 {
		res.Errors = append(res.Errors, Error{Pointer: pointerFor(ve), Message: validationMessage(ve)})
		return
	}
	for _, cause := range ve.Causes {
		collectLeaves(cause, res)
	}
}

func validationMessage(ve *jsonschema.ValidationError) string {
	if ve == nil {
		return "validation failed"
	}
	if ve.ErrorKind != nil {
		return ve.ErrorKind.LocalizedString(schemaPrinter)
	}
	return "validation failed"
}

func pointerFor(ve *jsonschema.ValidationError) string {
	base := jsonPointer(ve.InstanceLocation)
	if prop, ok := propertyFromKind(ve.ErrorKind); ok {
		return joinPointer(base, prop)
	}
	return base
}

func jsonPointer(tokens []string) string {
	if len(tokens) == 0 {
		return ""
	}
	var b strings.Builder
	for _, tok := range tokens {
		b.WriteByte('/')
		b.WriteString(escapePointerToken(tok))
	}
	return b.String()
}

func propertyFromKind(k jsonschema.ErrorKind) (string, bool) {
	if k == nil {
		return "", false
	}
	switch typed := k.(type) {
	case *kind.Required:
		if len(typed.Missing) == 1 {
			return typed.Missing[0], true
		}
	case *kind.AdditionalProperties:
		if len(typed.Properties) == 1 {
			return typed.Properties[0], true
		}
	}
	return "", false
}

func joinPointer(base, prop string) string {
	escaped := escapePointerToken(prop)
	if base == "" {
		return "/" + escaped
	}
	return base + "/" + escaped
}

func escapePointerToken(s string) string {
	s = strings.ReplaceAll(s, "~", "~0")
	s = strings.ReplaceAll(s, "/", "~1")
	return s
}
