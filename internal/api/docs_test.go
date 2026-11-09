package api

import "testing"

// Every setting is described: the generated reference prints both.
func TestSettingDefsDescribed(t *testing.T) {
	for _, d := range settingDefs {
		if d.desc == "" || d.typ == "" {
			t.Errorf("setting %s has no type or description", d.key)
		}
	}
}
