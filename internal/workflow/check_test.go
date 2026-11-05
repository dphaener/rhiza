package workflow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	for _, n := range BuiltinNames() {
		data, _ := Builtin(n)
		if _, ps := Validate(n, data, nil); len(ps) > 0 && !strings.Contains(ps[0].Message, "cannot be resolved") {
			t.Errorf("built-in %s: %v", n, ps)
		}
	}
	data, err := os.ReadFile(filepath.Join("testdata", "bad-outcome.toml"))
	if err != nil {
		t.Fatal(err)
	}
	_, ps := Validate("x", data, nil)
	if len(ps) != 1 || ps[0].Step != "build" || !strings.Contains(ps[0].Message, `on.sucess`) {
		t.Errorf("problems = %+v", ps)
	}
}
