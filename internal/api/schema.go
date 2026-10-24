package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/dphaener/genesis/internal/id"
)

// idPattern matches the canonical spelling of an id.ID.
const idPattern = "^[0-9A-HJKMNP-TV-Z]{26}$"

var (
	idType      = reflect.TypeFor[id.ID]()
	rawType     = reflect.TypeFor[json.RawMessage]()
	rawMapType  = reflect.TypeFor[map[string]json.RawMessage]()
	textMapType = reflect.TypeFor[map[string]string]()
)

// Schema returns the JSON Schema (a draft 2020-12 subset) of a request or
// response type, ready for canon.Marshal.
//
// A struct is an object with additionalProperties false. Every exported
// field must have a json tag: its name is the property name, "-" skips it,
// and a field without omitempty or omitzero is required. Unexported fields
// are skipped. A desc tag becomes the property's description; an enum tag
// ("a,b") lists its allowed strings, on the items of a slice.
//
// Strings (named string types too) are strings, bools booleans and ints
// integers. Slices are arrays, pointers are their element, id.ID is a
// pattern-checked string, json.RawMessage is any value ({}),
// map[string]json.RawMessage an object and map[string]string an object of
// strings. Any other type, an untagged exported field, an embedded field or
// an enum on a non-string panics: schemas are fixed by the code, so these are
// programming errors. Recursive types are not supported.
func Schema(t reflect.Type) map[string]any {
	return schemaOf(t, nil)
}

func schemaOf(t reflect.Type, enum []string) map[string]any {
	switch t {
	case idType:
		return withEnum(t, map[string]any{"type": "string", "pattern": idPattern}, enum)
	case rawType:
		return withEnum(t, map[string]any{}, enum)
	case rawMapType:
		return withEnum(t, map[string]any{"type": "object"}, enum)
	case textMapType:
		return withEnum(t, map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "string"}}, enum)
	}
	switch t.Kind() {
	case reflect.String:
		s := map[string]any{"type": "string"}
		if enum != nil {
			s["enum"] = enum
		}
		return s
	case reflect.Bool:
		return withEnum(t, map[string]any{"type": "boolean"}, enum)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return withEnum(t, map[string]any{"type": "integer"}, enum)
	case reflect.Slice:
		return map[string]any{"type": "array", "items": schemaOf(t.Elem(), enum)}
	case reflect.Pointer:
		return schemaOf(t.Elem(), enum)
	case reflect.Struct:
		return withEnum(t, structSchema(t), enum)
	}
	panic(fmt.Sprintf("api: schema: unsupported type %s", t))
}

// withEnum returns s, panicking if an enum was asked of a non-string type.
func withEnum(t reflect.Type, s map[string]any, enum []string) map[string]any {
	if enum != nil {
		panic(fmt.Sprintf("api: schema: enum tag on non-string type %s", t))
	}
	return s
}

func structSchema(t reflect.Type) map[string]any {
	props := map[string]any{}
	required := []string{}
	for i := range t.NumField() {
		f := t.Field(i)
		if f.Anonymous {
			panic(fmt.Sprintf("api: schema: %s.%s: embedded fields are not supported", t, f.Name))
		}
		if !f.IsExported() {
			continue
		}
		tag, ok := f.Tag.Lookup("json")
		if !ok {
			panic(fmt.Sprintf("api: schema: %s.%s: no json tag", t, f.Name))
		}
		if tag == "-" {
			continue
		}
		name, opts, _ := strings.Cut(tag, ",")
		if name == "" {
			panic(fmt.Sprintf("api: schema: %s.%s: json tag has no name", t, f.Name))
		}
		var enum []string
		if e, ok := f.Tag.Lookup("enum"); ok {
			enum = strings.Split(e, ",")
		}
		s := schemaOf(f.Type, enum)
		if d, ok := f.Tag.Lookup("desc"); ok {
			s["description"] = d
		}
		props[name] = s
		optional := slices.ContainsFunc(strings.Split(opts, ","), func(o string) bool {
			return o == "omitempty" || o == "omitzero"
		})
		if !optional {
			required = append(required, name)
		}
	}
	slices.Sort(required)
	return map[string]any{
		"type":                 "object",
		"properties":           props,
		"required":             required,
		"additionalProperties": false,
	}
}
