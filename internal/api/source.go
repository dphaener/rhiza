package api

// This file holds source add: a node of kind source made from a JSON file
// of the source's distillate (title, summary, quotes) and the sha256 of its
// raw material. The raw material itself, and its path, never reach the log.

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/canon"
)

func init() {
	register(Operation{
		Name:     "source.add",
		Summary:  "Add a source node from a JSON file of its title, summary and quotes, with the sha256 of its raw material",
		Write:    true,
		Record:   true,
		Request:  reflect.TypeFor[SourceAddRequest](),
		Response: reflect.TypeFor[WriteResponse](),
		Input: &Input{
			Summary: `{"title":"...","summary":"...","quotes":[{"text":"...","by":"...","where":"..."}]}: title, summary and quotes required (quotes may be empty); a quote needs text, by and where are optional`,
			Schema:  sourceFileSchema,
		},
		Run: runSourceAdd,
	})
}

// SourceAddRequest adds a source node.
type SourceAddRequest struct {
	Handle     string `json:"handle" arg:"0" desc:"The source's handle: lower-case words joined by single dashes, unique among live nodes"`
	File       string `json:"file" arg:"1" desc:"A JSON file of the source's title, summary and quotes (see the input schema)"`
	Raw        string `json:"raw,omitempty" desc:"The raw material's file, hashed with sha256; its content and path are not stored. Give one of --raw, --raw-sha256 and --raw-session"`
	RawSHA256  string `json:"raw_sha256,omitempty" desc:"The raw material's sha256, 64 lower-case hex characters. Give one of --raw, --raw-sha256 and --raw-session"`
	RawSession string `json:"raw_session,omitempty" desc:"A Rhiza session (full ID or unique prefix) whose agent transcript is the raw material, hashed with sha256 and named in a session field; its content and path are not stored. Give one of --raw, --raw-sha256 and --raw-session"`
}

// sourceFile is the decoded JSON file. Quotes is a pointer so a missing
// "quotes" can be told from an empty one.
type sourceFile struct {
	Title   string         `json:"title"`
	Summary string         `json:"summary"`
	Quotes  *[]sourceQuote `json:"quotes"`
}

// sourceQuote is one quote of a source.
type sourceQuote struct {
	Text  string `json:"text"`
	By    string `json:"by,omitempty"`
	Where string `json:"where,omitempty"`
}

// sourceFileSchema returns the JSON Schema of a source file.
func sourceFileSchema() map[string]any {
	str := func(desc string) map[string]any { return map[string]any{"type": "string", "description": desc} }
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"title":   map[string]any{"type": "string", "minLength": 1, "description": "The source's title"},
			"summary": map[string]any{"type": "string", "minLength": 1, "description": "What the source says, distilled"},
			"quotes": map[string]any{
				"type":        "array",
				"description": "Verbatim quotes from the source; may be empty",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"text":  map[string]any{"type": "string", "minLength": 1, "description": "The quote, verbatim"},
						"by":    str("Who said it"),
						"where": str("Where in the source it is"),
					},
					"required":             []string{"text"},
					"additionalProperties": false,
				},
			},
		},
		"required":             []string{"quotes", "summary", "title"},
		"additionalProperties": false,
	}
}

// runSourceAdd checks the flags, reads and decodes the file, hashes the raw
// material and adds the node through the add builder and write pipeline.
func runSourceAdd(ctx *Ctx, req any) (any, error) {
	r := req.(SourceAddRequest)
	const next = "rhiza describe source.add"
	given := 0
	for _, v := range []string{r.Raw, r.RawSHA256, r.RawSession} {
		if v != "" {
			given++
		}
	}
	if given != 1 {
		return nil, Usagef("source add needs exactly one of --raw, --raw-sha256 or --raw-session").WithNext(next)
	}
	if r.File == "" {
		return nil, Usagef("source add needs a file").WithNext(next)
	}
	data, err := readSourceFile("source file", r.File)
	if err != nil {
		return nil, err
	}
	f, err := decodeSourceFile(data)
	if err != nil {
		return nil, err
	}
	sum := r.RawSHA256
	var extra []string
	switch {
	case r.Raw != "":
		raw, err := readSourceFile("--raw", r.Raw)
		if err != nil {
			return nil, err
		}
		h := sha256.Sum256(raw)
		sum = hex.EncodeToString(h[:])
	case r.RawSession != "":
		sid, raw, err := sessionTranscript(ctx, r.RawSession)
		if err != nil {
			return nil, err
		}
		h := sha256.Sum256(raw)
		sum = hex.EncodeToString(h[:])
		extra = append(extra, "session="+sid)
	case !isSHA256Hex(sum):
		return nil, Usagef("--raw-sha256 %q: want 64 lower-case hex characters", sum).WithNext(next)
	}
	quotes, err := canon.Marshal(*f.Quotes)
	if err != nil {
		return nil, Errorf("quotes: %v", err)
	}
	add := AddRequest{
		Kind:       "source",
		Handle:     r.Handle,
		Fields:     append([]string{"title=" + f.Title, "summary=" + f.Summary, "raw_sha256=" + sum}, extra...),
		JSONFields: []string{"quotes=" + strings.TrimSuffix(string(quotes), "\n")},
	}
	return runWrite(ctx, "source.add", commandLine("source", "add", r.Handle),
		func(b writeBase) (edit, error) { return buildAdd(b, add) })
}

// sessionTranscript reads the agent transcript of the session ref names (a
// full ID or unique prefix) and returns the session's full ID with it. An
// unknown session is not_found; a session without a transcript is an error.
func sessionTranscript(ctx *Ctx, ref string) (string, []byte, error) {
	st, err := sessionStore(ctx)
	if err != nil {
		return "", nil, err
	}
	rec, err := readSession(st, strings.TrimSpace(ref))
	if err != nil {
		return "", nil, err
	}
	a, ok := adapter.Lookup(rec.Agent)
	if !ok {
		return "", nil, Errorf("session %s names agent %q, which this rhiza does not know: want one of %s",
			rec.ID, rec.Agent, strings.Join(adapter.Names(), ", "))
	}
	raw, err := readTranscript(ctx, a, rec.ID)
	if err != nil {
		return "", nil, err
	}
	return rec.ID, raw, nil
}

// readSourceFile reads file; a missing one is not_found, any other failure
// a usage error. what names the file in messages.
func readSourceFile(what, file string) ([]byte, error) {
	data, err := os.ReadFile(file)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return nil, NotFoundf("%s: no file %s", what, file)
	case err != nil:
		return nil, Usagef("%s: cannot read %s: %v", what, file, errors.Unwrap(err))
	}
	return data, nil
}

// decodeSourceFile decodes the file strictly: one JSON object with no
// duplicate or unknown keys, a non-empty title and summary, a quotes array
// and a non-empty text on every quote.
func decodeSourceFile(data []byte) (sourceFile, error) {
	const next = "rhiza describe source.add"
	trimmed := bytes.TrimSpace(data)
	if _, err := canon.Canonicalize(trimmed); err != nil {
		return sourceFile{}, Usagef("source file: %v", err).WithNext(next)
	}
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return sourceFile{}, Usagef("source file: not a JSON object").WithNext(next)
	}
	var f sourceFile
	dec := json.NewDecoder(bytes.NewReader(trimmed))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&f); err != nil {
		return sourceFile{}, Usagef("source file: %v", err).WithNext(next)
	}
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return sourceFile{}, Usagef("source file: trailing data").WithNext(next)
	}
	switch {
	case strings.TrimSpace(f.Title) == "":
		return sourceFile{}, Usagef(`source file: "title" is missing or empty`).WithNext(next)
	case strings.TrimSpace(f.Summary) == "":
		return sourceFile{}, Usagef(`source file: "summary" is missing or empty`).WithNext(next)
	case f.Quotes == nil:
		return sourceFile{}, Usagef(`source file: "quotes" is missing; give [] for none`).WithNext(next)
	}
	for i, q := range *f.Quotes {
		if strings.TrimSpace(q.Text) == "" {
			return sourceFile{}, Usagef(`source file: quotes[%d]: "text" is missing or empty`, i).WithNext(next)
		}
	}
	return f, nil
}

// isSHA256Hex reports whether s is 64 lower-case hex characters.
func isSHA256Hex(s string) bool {
	if len(s) != 64 {
		return false
	}
	for _, c := range s {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}
