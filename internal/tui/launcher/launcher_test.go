package launcher

import "testing"

func TestAppsAreSortedAndReplaced(t *testing.T) {
	mu.Lock()
	saved := registry
	registry = map[string]App{}
	mu.Unlock()
	defer func() { mu.Lock(); registry = saved; mu.Unlock() }()

	Register(App{Name: "show", Summary: "old"})
	Register(App{Name: "ready"})
	Register(App{Name: "show", Summary: "new"})
	got := Apps()
	if len(got) != 2 || got[0].Name != "ready" || got[1].Name != "show" || got[1].Summary != "new" {
		t.Fatalf("Apps() = %+v", got)
	}
}
