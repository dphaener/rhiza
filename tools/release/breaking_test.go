package main

import (
	"fmt"
	"strings"
	"testing"
)

func desc(ops ...string) []byte {
	return []byte(`{"operations":[` + strings.Join(ops, ",") + `]}`)
}

func op(name, req, resp string) string {
	return fmt.Sprintf(`{"name":%q,"request":{"type":"object","properties":{%s}},"response":{"type":"object","properties":{%s}}}`, name, req, resp)
}

func surf(d []byte, settings map[string]string, log int) Surface {
	return Surface{Describe: d, Settings: settings, LogFormat: log}
}

func TestBreaking(t *testing.T) {
	base := surf(desc(
		op("a", `"x":{"type":"string"},"y":{"type":"string"}`, `"r":{"type":"string"}`),
		op("b", ``, ``),
	), map[string]string{"s.one": "boolean", "s.two": "size"}, 5)
	tests := []struct {
		name string
		cur  Surface
		want []string
	}{
		{"identical", base, nil},
		{"removed op", surf(desc(op("a", `"x":{"type":"string"},"y":{"type":"string"}`, `"r":{"type":"string"}`)), base.Settings, 5),
			[]string{"operation: b removed"}},
		{"removed request field", surf(desc(op("a", `"x":{"type":"string"}`, `"r":{"type":"string"}`), op("b", ``, ``)), base.Settings, 5),
			[]string{"field: a request.y removed"}},
		{"removed response field", surf(desc(op("a", `"x":{"type":"string"},"y":{"type":"string"}`, ``), op("b", ``, ``)), base.Settings, 5),
			[]string{"field: a response.r removed"}},
		{"retyped field", surf(desc(op("a", `"x":{"type":"integer"},"y":{"type":"string"}`, `"r":{"type":"string"}`), op("b", ``, ``)), base.Settings, 5),
			[]string{`type: a request.x changed from "string" to "integer"`}},
		{"removed setting", surf(base.Describe, map[string]string{"s.one": "boolean"}, 5),
			[]string{"setting: s.two removed"}},
		{"retyped setting", surf(base.Describe, map[string]string{"s.one": "boolean", "s.two": "number"}, 5),
			[]string{`setting: s.two changed type from "size" to "number"`}},
		{"log format bump", surf(base.Describe, base.Settings, 6),
			[]string{"log-format: version changed from 5 to 6"}},
		{"additions are fine", surf(desc(
			op("a", `"x":{"type":"string"},"y":{"type":"string"},"z":{"type":"string"}`, `"r":{"type":"string"},"q":{"type":"string"}`),
			op("b", ``, ``), op("c", ``, ``)),
			map[string]string{"s.one": "boolean", "s.two": "size", "s.three": "boolean"}, 5), nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Breaking(base, tc.cur)
			if err != nil {
				t.Fatal(err)
			}
			var g []string
			for _, b := range got {
				g = append(g, b.String())
			}
			if strings.Join(g, "\n") != strings.Join(tc.want, "\n") {
				t.Fatalf("got %q want %q", g, tc.want)
			}
		})
	}
}

// writes is an apply-like input document: an array of write variants told apart by
// their op const, each with its own fields.
func writes(add, set, on string) string {
	v := []string{}
	if add != "" {
		v = append(v, `{"additionalProperties":false,"properties":{"op":{"const":"add"},`+add+`}}`)
	}
	if set != "" {
		v = append(v, `{"additionalProperties":false,"properties":{"op":{"const":"set"},`+set+`}}`)
	}
	return withInput(op("apply", "", `"on":{"type":"object","additionalProperties":{"type":"`+on+`"}}`),
		`"writes":{"type":"array","items":{"oneOf":[`+strings.Join(v, ",")+`]}}`)
}

// withInput adds an input document schema with props to an op.
func withInput(o, props string) string {
	return strings.TrimSuffix(o, "}") + `,"input":{"type":"object","properties":{` + props + `}}}`
}

func TestBreakingOneOfAndMaps(t *testing.T) {
	both := `"handle":{"type":"string"},"kind":{"type":"string"}`
	base := surf(desc(writes(both, `"handle":{"type":"string"}`, "string")), nil, 1)
	tests := []struct {
		name string
		cur  string
		want []string
	}{
		{"identical", writes(both, `"handle":{"type":"string"}`, "string"), nil},
		{"variants reordered", withInput(op("apply", "", `"on":{"type":"object","additionalProperties":{"type":"string"}}`),
			`"writes":{"type":"array","items":{"oneOf":[`+
				`{"properties":{"op":{"const":"set"},"handle":{"type":"string"}}},`+
				`{"properties":{"op":{"const":"add"},`+both+`}}]}}`), nil},
		{"field removed in a variant", writes(`"handle":{"type":"string"}`, `"handle":{"type":"string"}`, "string"),
			[]string{"field: apply input.writes[](op=add).kind removed"}},
		{"field retyped in a variant", writes(both, `"handle":{"type":"integer"}`, "string"),
			[]string{`type: apply input.writes[](op=set).handle changed from "string" to "integer"`}},
		{"variant dropped", writes(both, "", "string"),
			[]string{"variant: apply input.writes[](op=set) removed"}},
		{"variant added", withInput(op("apply", "", `"on":{"type":"object","additionalProperties":{"type":"string"}}`),
			`"writes":{"type":"array","items":{"oneOf":[`+
				`{"properties":{"op":{"const":"add"},`+both+`}},`+
				`{"properties":{"op":{"const":"set"},"handle":{"type":"string"}}},`+
				`{"properties":{"op":{"const":"mv"},"handle":{"type":"string"}}}]}}`), nil},
		{"map value retyped", writes(both, `"handle":{"type":"string"}`, "integer"),
			[]string{`type: apply response.on.* changed from "string" to "integer"`}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Breaking(base, surf(desc(tc.cur), nil, 1))
			if err != nil {
				t.Fatal(err)
			}
			var g []string
			for _, b := range got {
				g = append(g, b.String())
			}
			if strings.Join(g, "\n") != strings.Join(tc.want, "\n") {
				t.Fatalf("got %q want %q", g, tc.want)
			}
		})
	}
}

func TestParseSettingsAndFormat(t *testing.T) {
	md := "## M\n\n### `a.b`\n\nText.\n\n- Type: boolean\n- Default: `false`\n\n### `c.d`\n\n- Type: off, main or latest\n"
	got := parseSettings(md)
	if got["a.b"] != "boolean" || got["c.d"] != "off, main or latest" || len(got) != 2 {
		t.Fatalf("%v", got)
	}
	n, err := parseLogFormat("package x\n\nconst FormatVersion = 7\n")
	if err != nil || n != 7 {
		t.Fatal(n, err)
	}
}
