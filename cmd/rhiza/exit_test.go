package main

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/api"
)

// exitResponse is a test response that sets the exit code.
type exitResponse struct {
	Code int `json:"code"`
}

func (r exitResponse) Human(w io.Writer) { fmt.Fprintf(w, "exit %d\n", r.Code) }
func (r exitResponse) ExitCode() int     { return r.Code }

type exitRequest struct {
	Code int `json:"code"`
}

// exitOp is built by hand, not registered, so describe and help are
// unchanged.
var exitOp = api.Operation{
	Name:     "test.exit",
	Request:  reflect.TypeFor[exitRequest](),
	Response: reflect.TypeFor[exitResponse](),
	Run: func(_ *api.Ctx, req any) (any, error) {
		return exitResponse{Code: req.(exitRequest).Code}, nil
	},
}

func TestCallExiter(t *testing.T) {
	t.Chdir(t.TempDir())
	getenv := func(string) string { return "" }
	for _, code := range []int{0, 1} {
		var out [3]string
		for i, g := range []globals{{}, {json: true}, {}} {
			req := any(exitRequest{Code: code})
			if i == 2 {
				// The api call path: a decoded JSON request, JSON output.
				g.json = true
				var err error
				req, err = api.DecodeRequest(exitOp, fmt.Appendf(nil, `{"code":%d}`, code))
				if err != nil {
					t.Fatal(err)
				}
			}
			var stdout, stderr bytes.Buffer
			got := call(g, exitOp, req, strings.NewReader(""), &stdout, &stderr, getenv)
			if got != code {
				t.Errorf("code %d, path %d: exit %d, want %d", code, i, got, code)
			}
			if stdout.Len() == 0 {
				t.Errorf("code %d, path %d: empty stdout", code, i)
			}
			if stderr.Len() != 0 {
				t.Errorf("code %d, path %d: stderr %q", code, i, stderr.String())
			}
			out[i] = stdout.String()
		}
		if want := fmt.Sprintf("exit %d\n", code); out[0] != want {
			t.Errorf("code %d: human %q, want %q", code, out[0], want)
		}
		if want := fmt.Sprintf("{\n  \"code\": %d\n}\n", code); out[1] != want {
			t.Errorf("code %d: json %q, want %q", code, out[1], want)
		}
		if out[1] != out[2] {
			t.Errorf("code %d: --json %q differs from api call %q", code, out[1], out[2])
		}
	}
}

// TestCallExitNoExiter checks a response without ExitCode still exits 0.
func TestCallExitNoExiter(t *testing.T) {
	t.Chdir(t.TempDir())
	var stdout, stderr bytes.Buffer
	if got := call(globals{json: true}, testOp, testRequest{}, strings.NewReader(""), &stdout, &stderr, func(string) string { return "" }); got != 0 {
		t.Errorf("exit %d, want 0", got)
	}
}
