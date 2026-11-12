package main

import "testing"

func TestVersionNewer(t *testing.T) {
	cases := []struct {
		newest, current string
		want            bool
	}{
		{"v0.2.0", "v0.1.9", true},
		{"v0.1.9", "v0.2.0", false},
		{"v0.2.0", "v0.2.0", false},
		{"v0.2.0", "v0.2.0-rc.1", true},
		{"v0.2.0-rc.1", "v0.2.0", false},
		{"v0.2.0-rc.10", "v0.2.0-rc.2", true},
		{"v0.2.0-rc.2", "v0.2.0-rc.10", false},
		{"v0.2.0-rc.1", "v0.2.0-rc", true},
		{"v0.2.0-rc", "v0.2.0-1", true},
		{"v0.2.1-main.10+abc", "v0.2.1-main.9+def", true},
		{"v0.2.1-main.9+def", "v0.2.1-main.10+abc", false},
		{"v0.2.0-5-gabc1234", "v0.2.0", true},
		{"v0.2.0", "v0.2.0-5-gabc1234", false},
		{"v0.2.1-12-gdef", "v0.2.1-5-gabc", true},
		{"v0.2.0-5-gabc1234", "v0.2.0-12-gdef5678", false},
		{"v0.2.0-12-gdef5678-dirty", "v0.2.0-5-gabc1234", true},
		{"v0.2.0-rc.1-3-gabc", "v0.2.0-rc.1", true},
		{"v0.2.0-rc.1-3-gabc", "v0.2.0", false},
		{"v0.2.1", "v0.2.0-40-gabc", true},
		{"v0.2.0", "dev", true},
		{"dev", "dev", false},
	}
	for _, c := range cases {
		if got := versionNewer(c.newest, c.current); got != c.want {
			t.Errorf("versionNewer(%q, %q) = %v, want %v", c.newest, c.current, got, c.want)
		}
	}
}
