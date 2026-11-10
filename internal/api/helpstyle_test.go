package api

import (
	"reflect"
	"strings"
	"testing"
	"unicode"
)

// TestHelpStyle holds every operation summary and every desc tag on a
// request struct to the CLI help style: a capitalised sentence with no
// trailing period, and defaults written as "; default <x>".
func TestHelpStyle(t *testing.T) {
	for _, o := range Operations() {
		check := func(where, s string) {
			if s == "" {
				return
			}
			if r := []rune(s)[0]; !unicode.IsUpper(r) && !unicode.IsDigit(r) {
				t.Errorf("%s: %q does not start with a capital letter", where, s)
			}
			if strings.HasSuffix(s, ".") {
				t.Errorf("%s: %q ends with a period", where, s)
			}
			low := strings.ToLower(s)
			for _, bad := range []string{"(default", "default:", ", default", ". default", "defaults to", "default is", "(defaults"} {
				if strings.Contains(low, bad) {
					t.Errorf("%s: %q writes its default as %q, not \"; default <x>\"", where, s, bad)
				}
			}
		}
		check(o.Name+" summary", o.Summary)
		walkRequest(o.Request, map[reflect.Type]bool{}, func(f reflect.StructField) {
			if d, ok := f.Tag.Lookup("desc"); ok {
				check(o.Name+"."+f.Name, d)
			}
		})
	}
}

func walkRequest(rt reflect.Type, seen map[reflect.Type]bool, fn func(reflect.StructField)) {
	for rt.Kind() == reflect.Pointer || rt.Kind() == reflect.Slice || rt.Kind() == reflect.Map {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct || seen[rt] {
		return
	}
	seen[rt] = true
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		fn(f)
		walkRequest(f.Type, seen, fn)
	}
}
