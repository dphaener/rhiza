package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/daemon"
)

// The event command's payload carries a change set's ref verbatim on every
// move, and has no ref key for a change set without one.
func TestDaemonEventsRef(t *testing.T) {
	ctx := testCtx(t.TempDir())
	ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
	call := func(name, req string) any {
		t.Helper()
		resp, err := callOp(t, ctx, name, req)
		if err != nil {
			t.Fatalf("%s %s: %v", name, req, err)
		}
		return resp
	}
	call("init", `{}`)
	out := filepath.Join(t.TempDir(), "payloads")
	cmd := `cat >> ` + out + `; echo >> ` + out
	var state daemon.State
	if err := daemonEvents(ctx, &state, cmd); err != nil { // seeds
		t.Fatal(err)
	}
	const ref = "https://example.com/issues/7?x=\"a b\"&y=é"
	withRef := call("change.new", `{"intent":"With ref","ref":"`+strings.ReplaceAll(ref, `"`, `\"`)+`"}`).(ChangeNewResponse).ChangeSet.String()
	call("propose", `{}`)
	call("accept", `{"reason":"ok"}`)
	without := call("change.new", `{"intent":"Without ref"}`).(ChangeNewResponse).ChangeSet.String()
	call("propose", `{}`)
	if err := daemonEvents(ctx, &state, cmd); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]int{}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if line == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("payload %q: %v", line, err)
		}
		cs, _ := m["change"].(string)
		got, has := m["ref"]
		switch cs {
		case withRef:
			if got != ref {
				t.Errorf("%v payload ref = %v, want %q", m["event"], got, ref)
			}
		case without:
			if has {
				t.Errorf("%v payload has a ref key: %v", m["event"], got)
			}
		}
		seen[cs]++
	}
	if seen[withRef] != 2 || seen[without] != 1 {
		t.Errorf("payloads per change set = %v, want 2 with ref and 1 without", seen)
	}
}
