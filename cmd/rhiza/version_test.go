package main

import (
	"os"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	for _, tc := range []struct{ name, ldflags, module, want string }{
		{"release build stamps the tag", "v0.1.0", "(devel)", "v0.1.0"},
		{"ldflags win over the module", "v0.1.0", "v0.2.0", "v0.1.0"},
		{"go install reports the module version", "", "v0.1.0", "v0.1.0"},
		{"local build", "", "(devel)", "dev"},
		{"no build info", "", "", "dev"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := pickVersion(tc.ldflags, tc.module); got != tc.want {
				t.Fatalf("pickVersion(%q, %q) = %q, want %q", tc.ldflags, tc.module, got, tc.want)
			}
		})
	}
}

// TestVersionGoReleaserStamp keeps the release config stamping the tag
// (with its v), not {{ .Version }}, which GoReleaser strips of the v.
func TestVersionGoReleaserStamp(t *testing.T) {
	b, err := os.ReadFile("../../.goreleaser.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "-X main.version={{ .Tag }}") {
		t.Fatal(".goreleaser.yaml must set -X main.version={{ .Tag }}")
	}
}
