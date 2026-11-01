package changeset

import (
	"reflect"
	"testing"
)

func TestLineDiff(t *testing.T) {
	cases := []struct {
		name, old, new string
		want           []Line
	}{
		{"middle changed", "one\ntwo\nthree", "one\nTWO\nthree", []Line{
			{LineSame, "one"}, {LineRemoved, "two"}, {LineAdded, "TWO"}, {LineSame, "three"}}},
		{"same", "a\nb", "a\nb", []Line{{LineSame, "a"}, {LineSame, "b"}}},
		{"line added", "a\nc", "a\nb\nc", []Line{{LineSame, "a"}, {LineAdded, "b"}, {LineSame, "c"}}},
		{"line removed", "a\nb\nc", "a\nc", []Line{{LineSame, "a"}, {LineRemoved, "b"}, {LineSame, "c"}}},
		{"single to multi", "a", "a\nb", []Line{{LineSame, "a"}, {LineAdded, "b"}}},
	}
	for _, c := range cases {
		if got := LineDiff(c.old, c.new); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}
