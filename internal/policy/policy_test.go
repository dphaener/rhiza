package policy

import (
	"strings"
	"testing"
	"testing/fstest"
)

func TestRead(t *testing.T) {
	file := func(s string) fstest.MapFS {
		return fstest.MapFS{Path: {Data: []byte(s)}}
	}
	tests := []struct {
		name    string
		fs      fstest.MapFS
		want    string
		wantErr []string
	}{
		{"missing file", fstest.MapFS{}, BacklogAfter, nil},
		{"empty file", file(""), BacklogAfter, nil},
		{"missing key", file("[schedule]\n"), BacklogAfter, nil},
		{"other table", file("[other]\nbacklog = \"skip\"\n"), BacklogAfter, nil},
		{"after", file("[schedule]\nbacklog = \"after\"\n"), BacklogAfter, nil},
		{"before", file("[schedule]\nbacklog = \"before\"\n"), BacklogBefore, nil},
		{"skip", file("[schedule]\nbacklog = \"skip\"\n"), BacklogSkip, nil},
		{"unknown value", file("[schedule]\nbacklog = \"sometimes\"\n"), BacklogAfter,
			[]string{".rhiza/policy.toml", `"sometimes"`}},
		{"empty value", file("[schedule]\nbacklog = \"\"\n"), BacklogAfter, []string{`""`}},
		{"wrong type", file("[schedule]\nbacklog = 3\n"), BacklogAfter, []string{".rhiza/policy.toml"}},
		{"malformed", file("[schedule\nbacklog = "), BacklogAfter, []string{".rhiza/policy.toml"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(tt.fs)
			if got.Backlog != tt.want {
				t.Errorf("Backlog = %q, want %q", got.Backlog, tt.want)
			}
			if len(tt.wantErr) == 0 {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("no error, want one naming %q", tt.wantErr)
			}
			for _, s := range tt.wantErr {
				if !strings.Contains(err.Error(), s) {
					t.Errorf("error %q does not contain %q", err, s)
				}
			}
		})
	}
}
