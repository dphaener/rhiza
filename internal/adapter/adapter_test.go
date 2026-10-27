package adapter

import (
	"reflect"
	"testing"
)

type fake struct{ name string }

func (f fake) Name() string            { return f.name }
func (fake) PermissionModes() []string { return []string{"safe", "risky"} }
func (fake) Bypass(mode string) bool   { return mode == "risky" }
func (f fake) Command(l Launch) ([]string, error) {
	if err := CheckMode(f, l.PermissionMode); err != nil {
		return nil, err
	}
	return []string{f.name, l.PermissionMode}, nil
}
func (fake) HookStatus(string, []byte) (HookReport, bool, error)   { return HookReport{}, false, nil }
func (fake) AskReply(string, []byte, string) ([]byte, bool, error) { return nil, false, nil }
func (fake) HookMessage(string, []byte) (Message, bool, error)     { return Message{}, false, nil }
func (fake) TranscriptGlobs(string, func(string) string) []string  { return nil }
func (fake) Answers([]byte) ([]Answer, error)                      { return nil, nil }
func (fake) TokenUse([]byte) (TokenUse, error)                     { return TokenUse{}, nil }

func reset() {
	mu.Lock()
	adapters = map[string]Adapter{}
	def = ""
	mu.Unlock()
}

func TestRegistry(t *testing.T) {
	reset()
	t.Cleanup(reset)
	if Default() != nil {
		t.Fatal("Default with an empty registry should be nil")
	}
	Register(fake{"b"})
	RegisterDefault(fake{"a"})
	if got := Names(); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("Names = %v", got)
	}
	if a, ok := Lookup("b"); !ok || a.Name() != "b" {
		t.Fatalf("Lookup(b) = %v, %v", a, ok)
	}
	if _, ok := Lookup("zzz"); ok {
		t.Fatal("Lookup of an unknown name should fail")
	}
	if Default().Name() != "a" {
		t.Fatalf("Default = %v", Default())
	}
}

func TestRegisterDuplicatePanics(t *testing.T) {
	reset()
	t.Cleanup(reset)
	Register(fake{"a"})
	defer func() {
		if recover() == nil {
			t.Fatal("duplicate Register should panic")
		}
	}()
	Register(fake{"a"})
}

func TestCheckMode(t *testing.T) {
	f := fake{"a"}
	for _, m := range []string{"", "yolo"} {
		if argv, err := f.Command(Launch{PermissionMode: m}); err == nil || argv != nil {
			t.Errorf("mode %q: argv=%v err=%v", m, argv, err)
		}
	}
	if _, err := f.Command(Launch{PermissionMode: "safe"}); err != nil {
		t.Fatal(err)
	}
}
