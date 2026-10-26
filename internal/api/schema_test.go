package api

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

type schemaInner struct {
	Note string `json:"note,omitempty" desc:"A note."`
	N    int    `json:"n"`
}

type schemaSample struct {
	Kind     op.Kind                    `json:"kind" desc:"The node kind." enum:"intent,capability"`
	Handle   string                     `json:"handle" desc:"The handle."`
	Tags     []string                   `json:"tags,omitempty" desc:"Tags."`
	Modes    []string                   `json:"modes,omitempty" enum:"a,b"`
	Parent   id.ID                      `json:"parent,omitempty" desc:"The parent."`
	Done     bool                       `json:"done"`
	Inner    schemaInner                `json:"inner"`
	Ptr      *schemaInner               `json:"ptr,omitempty"`
	Items    []schemaInner              `json:"items"`
	Raw      json.RawMessage            `json:"raw,omitempty"`
	Fields   map[string]json.RawMessage `json:"fields,omitempty"`
	Labels   map[string]string          `json:"labels,omitempty"`
	Skipped  string                     `json:"-"`
	internal string
}

const schemaSampleWant = `{
  "additionalProperties": false,
  "properties": {
    "done": {
      "type": "boolean"
    },
    "fields": {
      "type": "object"
    },
    "handle": {
      "description": "The handle.",
      "type": "string"
    },
    "inner": {
      "additionalProperties": false,
      "properties": {
        "n": {
          "type": "integer"
        },
        "note": {
          "description": "A note.",
          "type": "string"
        }
      },
      "required": [
        "n"
      ],
      "type": "object"
    },
    "items": {
      "items": {
        "additionalProperties": false,
        "properties": {
          "n": {
            "type": "integer"
          },
          "note": {
            "description": "A note.",
            "type": "string"
          }
        },
        "required": [
          "n"
        ],
        "type": "object"
      },
      "type": "array"
    },
    "kind": {
      "description": "The node kind.",
      "enum": [
        "intent",
        "capability"
      ],
      "type": "string"
    },
    "labels": {
      "additionalProperties": {
        "type": "string"
      },
      "type": "object"
    },
    "modes": {
      "items": {
        "enum": [
          "a",
          "b"
        ],
        "type": "string"
      },
      "type": "array"
    },
    "parent": {
      "description": "The parent.",
      "pattern": "^[0-9A-HJKMNP-TV-Z]{26}$",
      "type": "string"
    },
    "ptr": {
      "additionalProperties": false,
      "properties": {
        "n": {
          "type": "integer"
        },
        "note": {
          "description": "A note.",
          "type": "string"
        }
      },
      "required": [
        "n"
      ],
      "type": "object"
    },
    "raw": {},
    "tags": {
      "description": "Tags.",
      "items": {
        "type": "string"
      },
      "type": "array"
    }
  },
  "required": [
    "done",
    "handle",
    "inner",
    "items",
    "kind"
  ],
  "type": "object"
}
`

func TestSchema(t *testing.T) {
	b, err := canon.Marshal(Schema(reflect.TypeFor[schemaSample]()))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != schemaSampleWant {
		t.Errorf("got\n%s\nwant\n%s", b, schemaSampleWant)
	}
}

func TestSchemaEmptyStruct(t *testing.T) {
	b, err := canon.Marshal(Schema(reflect.TypeFor[struct{}]()))
	if err != nil {
		t.Fatal(err)
	}
	want := "{\n  \"additionalProperties\": false,\n  \"properties\": {},\n  \"required\": [],\n  \"type\": \"object\"\n}\n"
	if string(b) != want {
		t.Errorf("got\n%s\nwant\n%s", b, want)
	}
}

func TestSchemaScalars(t *testing.T) {
	cases := []struct {
		t    reflect.Type
		want string
	}{
		{reflect.TypeFor[string](), `{"type":"string"}`},
		{reflect.TypeFor[bool](), `{"type":"boolean"}`},
		{reflect.TypeFor[int64](), `{"type":"integer"}`},
		{reflect.TypeFor[uint](), `{"type":"integer"}`},
		{reflect.TypeFor[id.ID](), `{"pattern":"^[0-9A-HJKMNP-TV-Z]{26}$","type":"string"}`},
		{reflect.TypeFor[[]id.ID](), `{"items":{"pattern":"^[0-9A-HJKMNP-TV-Z]{26}$","type":"string"},"type":"array"}`},
		{reflect.TypeFor[json.RawMessage](), `{}`},
	}
	for _, c := range cases {
		t.Run(c.t.String(), func(t *testing.T) {
			got, err := canon.Marshal(Schema(c.t))
			if err != nil {
				t.Fatal(err)
			}
			want, err := canon.Canonicalize([]byte(c.want))
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != string(want) {
				t.Errorf("got %s, want %s", got, want)
			}
		})
	}
}

func TestSchemaPanics(t *testing.T) {
	type untagged struct{ Name string }
	type float struct {
		X float64 `json:"x"`
	}
	type badMap struct {
		M map[string]int `json:"m"`
	}
	type badEnum struct {
		B bool `json:"b" enum:"x,y"`
	}
	for name, typ := range map[string]reflect.Type{
		"untagged": reflect.TypeFor[untagged](),
		"float":    reflect.TypeFor[float](),
		"badMap":   reflect.TypeFor[badMap](),
		"badEnum":  reflect.TypeFor[badEnum](),
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Error("no panic")
				}
			}()
			Schema(typ)
		})
	}
}
