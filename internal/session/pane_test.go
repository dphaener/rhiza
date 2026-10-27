package session

import "testing"

func TestSplitFor(t *testing.T) {
	for _, c := range []struct {
		name     string
		w, h     int
		override string
		want     string
	}{
		{"landscape", 1600, 900, "auto", SplitRight},
		{"portrait", 900, 1600, "auto", SplitBelow},
		{"square", 1000, 1000, "", SplitRight},
		{"portrait no override", 900, 1600, "", SplitBelow},
		{"override below on landscape", 1600, 900, SplitBelow, SplitBelow},
		{"override right on portrait", 900, 1600, SplitRight, SplitRight},
		{"unknown size", 0, 0, "auto", SplitRight},
		{"unknown height", 900, 0, "auto", SplitRight},
		{"unknown size with override", 0, 0, SplitBelow, SplitBelow},
		{"unknown override", 900, 1600, "sideways", SplitBelow},
	} {
		if got := SplitFor(c.w, c.h, c.override); got != c.want {
			t.Errorf("%s: SplitFor(%d, %d, %q) = %q, want %q", c.name, c.w, c.h, c.override, got, c.want)
		}
	}
}
