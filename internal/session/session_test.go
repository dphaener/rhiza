package session

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestRecordJSONFieldNames(t *testing.T) {
	r := Record{
		ID: "i", Agent: "a", Change: "c", Title: "t", Step: "build", Actor: "x", For: "y",
		PermissionMode: "manual", Worktree: "/w", Branch: "cs/t",
		Tab:       Tab{Session: "s", Window: "w", Block: "b", Title: "T"},
		StartedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"id"`, `"agent"`, `"change"`, `"title"`, `"step"`, `"actor"`, `"for"`,
		`"permission_mode"`, `"worktree"`, `"branch"`, `"tab"`, `"session"`, `"window"`, `"block"`,
		`"started_at":"2026-01-02T03:04:05Z"`, `"status"`} {
		if !strings.Contains(string(data), key) {
			t.Errorf("missing %s in %s", key, data)
		}
	}
	if strings.Contains(string(data), `"run"`) {
		t.Errorf("empty run should be omitted: %s", data)
	}
}

func TestNewID(t *testing.T) {
	cases := []struct {
		in   [16]byte
		want string
	}{
		{[16]byte{}, "00000000-0000-4000-8000-000000000000"},
		{[16]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, "ffffffff-ffff-4fff-bfff-ffffffffffff"},
		{[16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, "01020304-0506-4708-890a-0b0c0d0e0f10"},
	}
	for _, c := range cases {
		if got := NewID(c.in); got != c.want {
			t.Errorf("NewID(%v) = %s, want %s", c.in, got, c.want)
		}
	}
}
