package api

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

type paramsRequest struct {
	Kind    string   `json:"kind" arg:"0" desc:"The node kind" enum:"a,b"`
	Handle  string   `json:"handle,omitempty" arg:"1"`
	Parent  string   `json:"parent,omitempty" desc:"Parent ref"`
	Fields  []string `json:"fields,omitempty" flag:"field"`
	Force   bool     `json:"force_it,omitempty"`
	Depth   int      `json:"depth,omitempty"`
	Reason  string   `json:"reason"`
	Hidden  string   `json:"hidden,omitempty" flag:"-"`
	private string
}

func TestParams(t *testing.T) {
	ps := Params(reflect.TypeFor[paramsRequest]())
	type row struct {
		Name, Flag string
		Arg        int
		Kind       ParamKind
		Required   bool
	}
	var got []row
	for _, p := range ps {
		got = append(got, row{p.Name, p.Flag, p.Arg, p.Kind, p.Required})
	}
	want := []row{
		{"kind", "", 0, ParamString, true},
		{"handle", "", 1, ParamString, false},
		{"parent", "parent", -1, ParamString, false},
		{"fields", "field", -1, ParamList, false},
		{"force_it", "force-it", -1, ParamBool, false},
		{"depth", "depth", -1, ParamInt, false},
		{"reason", "reason", -1, ParamString, true},
	}
	if !slices.Equal(got, want) {
		t.Errorf("Params =\n%v\nwant\n%v", got, want)
	}
	if ps[0].Desc != "The node kind" || !slices.Equal(ps[0].Enum, []string{"a", "b"}) {
		t.Errorf("Params[0] desc/enum = %q %v", ps[0].Desc, ps[0].Enum)
	}
}

type restRequest struct {
	First string   `json:"first" arg:"0"`
	Rest  []string `json:"rest,omitempty" arg:"1"`
}

func TestParamsRest(t *testing.T) {
	ps := Params(reflect.TypeFor[restRequest]())
	if len(ps) != 2 || !ps[1].Rest || ps[1].Kind != ParamList || ps[0].Rest {
		t.Errorf("Params = %+v", ps)
	}
}

func TestParamsPanics(t *testing.T) {
	type gap struct {
		A string `json:"a" arg:"1"`
	}
	type restNotLast struct {
		A []string `json:"a" arg:"0"`
		B string   `json:"b" arg:"1"`
	}
	type badType struct {
		A map[string]string `json:"a"`
	}
	type badArg struct {
		A string `json:"a" arg:"x"`
	}
	cases := map[string]reflect.Type{
		"gap":           reflect.TypeFor[gap](),
		"rest not last": reflect.TypeFor[restNotLast](),
		"bad type":      reflect.TypeFor[badType](),
		"bad arg":       reflect.TypeFor[badArg](),
	}
	for name, typ := range cases {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Error("no panic")
				}
			}()
			Params(typ)
		})
	}
}

func TestUsageLine(t *testing.T) {
	o := Operation{Name: "change.new", Request: reflect.TypeFor[paramsRequest]()}
	got := UsageLine(o)
	want := "rhiza change new <kind> [<handle>] --reason <reason> [--parent <parent>] [--field <field>]... [--force-it] [--depth <depth>]"
	if got != want {
		t.Errorf("UsageLine =\n%s\nwant\n%s", got, want)
	}
	o.Request = reflect.TypeFor[restRequest]()
	if got := UsageLine(o); !strings.HasSuffix(got, "<first> [<rest>...]") {
		t.Errorf("UsageLine rest = %s", got)
	}
}
