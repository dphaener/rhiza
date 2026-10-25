package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"reflect"

	"github.com/dphaener/genesis/internal/canon"
)

// DecodeRequest decodes data into a new value of o.Request and returns that
// value itself, not a pointer to it, ready for o.Run. Empty or all-whitespace
// data means {}.
//
// Decoding is strict: data must be one JSON object with no duplicate keys,
// no unknown fields, no wrongly typed values, no non-integer numbers
// (Genesis JSON has only integers) and nothing after it. Every failure is a
// usage Error.
func DecodeRequest(o Operation, data []byte) (any, error) {
	ptr := reflect.New(o.Request)
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return ptr.Elem().Interface(), nil
	}
	// canon rejects malformed JSON, duplicate keys and trailing data, which
	// encoding/json would accept or report less clearly.
	if _, err := canon.Canonicalize(trimmed); err != nil {
		return nil, Usagef("%s request: %v", o.Name, err)
	}
	if trimmed[0] != '{' {
		return nil, Usagef("%s request: not a JSON object", o.Name)
	}
	dec := json.NewDecoder(bytes.NewReader(trimmed))
	dec.DisallowUnknownFields()
	if err := dec.Decode(ptr.Interface()); err != nil {
		return nil, Usagef("%s request: %v", o.Name, err)
	}
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return nil, Usagef("%s request: trailing data", o.Name)
	}
	return ptr.Elem().Interface(), nil
}
