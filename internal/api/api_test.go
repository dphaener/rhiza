package api

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

type emptyRequest struct{}
type emptyResponse struct{}

func testOp(name string) Operation {
	return Operation{
		Name:     name,
		Summary:  "Test " + name,
		Request:  reflect.TypeFor[emptyRequest](),
		Response: reflect.TypeFor[emptyResponse](),
		Run:      func(*Ctx, any) (any, error) { return emptyResponse{}, nil },
	}
}

func TestRegistrySorted(t *testing.T) {
	var r registry
	for _, n := range []string{"propose", "add", "change.new", "change.ls", "init"} {
		r.add(testOp(n))
	}
	var got []string
	for _, o := range r.all() {
		got = append(got, o.Name)
	}
	want := []string{"add", "change.ls", "change.new", "init", "propose"}
	if !slices.Equal(got, want) {
		t.Errorf("all() = %v, want %v", got, want)
	}
	o, ok := r.lookup("change.new")
	if !ok || o.Name != "change.new" {
		t.Errorf("lookup(change.new) = %v, %v", o.Name, ok)
	}
	if _, ok := r.lookup("nope"); ok {
		t.Error("lookup(nope) found something")
	}
}

func TestRegistryPanics(t *testing.T) {
	cases := []struct {
		name  string
		setup func(r *registry)
		op    Operation
		want  string
	}{
		{"duplicate", func(r *registry) { r.add(testOp("add")) }, testOp("add"), "duplicate"},
		{"empty name", func(*registry) {}, testOp(""), "empty"},
		{"no request", func(*registry) {}, func() Operation { o := testOp("x"); o.Request = nil; return o }(), "Request"},
		{"no response", func(*registry) {}, func() Operation { o := testOp("x"); o.Response = nil; return o }(), "Response"},
		{"no run", func(*registry) {}, func() Operation { o := testOp("x"); o.Run = nil; return o }(), "Run"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var r registry
			c.setup(&r)
			defer func() {
				p := recover()
				if p == nil {
					t.Fatal("no panic")
				}
				if !strings.Contains(p.(string), c.want) {
					t.Errorf("panic %q does not mention %q", p, c.want)
				}
			}()
			r.add(c.op)
		})
	}
}

func TestRegistryDefault(t *testing.T) {
	// No operations are registered yet; the package-level functions read
	// the default registry.
	if ops := Operations(); len(ops) != len(defaultRegistry.ops) {
		t.Errorf("Operations() has %d, registry %d", len(ops), len(defaultRegistry.ops))
	}
	if _, ok := Lookup("no.such.op"); ok {
		t.Error("Lookup found an unknown operation")
	}
}
