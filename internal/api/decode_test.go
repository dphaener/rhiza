package api

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/id"
)

type decodeRequest struct {
	Handle string   `json:"handle"`
	Tags   []string `json:"tags,omitempty"`
	Parent id.ID    `json:"parent,omitempty"`
}

func decodeOp() Operation {
	o := testOp("decode")
	o.Request = reflect.TypeFor[decodeRequest]()
	return o
}

func TestDecode(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want decodeRequest
	}{
		{"object", `{"handle":"core","tags":["a","b"]}`, decodeRequest{Handle: "core", Tags: []string{"a", "b"}}},
		{"id", `{"parent":"01JB0000000000000000000PAR"}`, decodeRequest{}},
		{"empty", ``, decodeRequest{}},
		{"whitespace", " \n\t ", decodeRequest{}},
		{"empty object", `{}`, decodeRequest{}},
		{"surrounding whitespace", "\n {\"handle\":\"x\"} \n", decodeRequest{Handle: "x"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := DecodeRequest(decodeOp(), []byte(c.in))
			if err != nil {
				t.Fatal(err)
			}
			req, ok := got.(decodeRequest)
			if !ok {
				t.Fatalf("got %T, want decodeRequest", got)
			}
			if c.name == "id" {
				// IDs are compared by their canonical spelling.
				if req.Parent.String() != "01JB0000000000000000000PAR" {
					t.Errorf("parent = %s", req.Parent)
				}
				return
			}
			if !reflect.DeepEqual(req, c.want) {
				t.Errorf("got %+v, want %+v", req, c.want)
			}
		})
	}
}

func TestDecodeRejects(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"unknown field", `{"bogus":1}`, "bogus"},
		{"trailing garbage", `{"handle":"x"} junk`, "trailing"},
		{"trailing object", `{"handle":"x"}{}`, "trailing"},
		{"wrong type", `{"handle":1}`, "handle"},
		{"bad id", `{"parent":"nope"}`, "invalid id"},
		{"not an object", `["x"]`, "object"},
		{"null", `null`, "object"},
		{"malformed", `{"handle":`, "request"},
		{"duplicate key", `{"handle":"a","handle":"b"}`, "duplicate"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := DecodeRequest(decodeOp(), []byte(c.in))
			if err == nil {
				t.Fatalf("got %+v, want an error", got)
			}
			if ExitCode(err) != 2 {
				t.Errorf("ExitCode = %d, want 2 (usage): %v", ExitCode(err), err)
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("error %q does not mention %q", err, c.want)
			}
		})
	}
}
