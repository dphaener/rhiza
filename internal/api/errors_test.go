package api

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dphaener/genesis/internal/canon"
)

func TestExitCode(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, 0},
		{"plain", errors.New("boom"), 1},
		{"error", Errorf("boom"), 1},
		{"usage", Usagef("bad %s", "flag"), 2},
		{"conflict", Conflictf("taken"), 3},
		{"not_found", NotFoundf("no %q", "x"), 4},
		{"wrapped usage", fmt.Errorf("outer: %w", Usagef("bad")), 2},
		{"wrapped not_found", fmt.Errorf("a: %w", fmt.Errorf("b: %w", NotFoundf("gone"))), 4},
		{"unknown code", &Error{Code: "weird", Message: "x"}, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := ExitCode(c.err); got != c.want {
				t.Errorf("ExitCode(%v) = %d, want %d", c.err, got, c.want)
			}
		})
	}
}

func TestErrorMessageAndNext(t *testing.T) {
	e := Conflictf("handle %q is taken", "core").WithNext("genesis show core")
	if e.Error() != `handle "core" is taken` {
		t.Errorf("Error() = %q", e.Error())
	}
	if e.Code != CodeConflict || e.Next != "genesis show core" {
		t.Errorf("got %+v", e)
	}
	base := Usagef("x")
	_ = base.WithNext("y")
	if base.Next != "" {
		t.Errorf("WithNext changed its receiver: %+v", base)
	}
}

func TestErrorJSON(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"api error", fmt.Errorf("wrap: %w", NotFoundf("no node %q", "x").WithNext("genesis add")),
			"{\n  \"error\": {\n    \"code\": \"not_found\",\n    \"message\": \"no node \\\"x\\\"\",\n    \"next\": \"genesis add\"\n  }\n}\n"},
		{"plain error", errors.New("disk full"),
			"{\n  \"error\": {\n    \"code\": \"error\",\n    \"message\": \"disk full\",\n    \"next\": \"\"\n  }\n}\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, err := canon.Marshal(ErrorJSON(c.err))
			if err != nil {
				t.Fatal(err)
			}
			if string(b) != c.want {
				t.Errorf("got\n%s\nwant\n%s", b, c.want)
			}
		})
	}
}

func TestArticle(t *testing.T) {
	for s, want := range map[string]string{"intent": "an", "evidence": "an", "actor": "an", "exempts": "an", "capability": "a", "decision": "a", "": "a"} {
		if got := article(s); got != want {
			t.Errorf("article(%q) = %q, want %q", s, got, want)
		}
	}
}
