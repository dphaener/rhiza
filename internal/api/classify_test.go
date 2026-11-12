package api

import (
	"errors"
	"testing"
)

func TestClassifyFailure(t *testing.T) {
	say := func(out string, err error) func(string, string) (string, error) {
		return func(string, string) (string, error) { return out, err }
	}
	for _, tc := range []struct {
		name      string
		command   string
		shell     func(string, string) (string, error)
		transient bool
		reason    string
	}{
		{"transient", "c", say("transient: GitHub 502\n", nil), true, "GitHub 502"},
		{"real", "c", say("real: assertion failed\n", nil), false, "assertion failed"},
		{"leading blank lines", "c", say("\n  Transient - lock held\n", nil), true, "lock held"},
		{"garbage", "c", say("I think maybe it is flaky?\n", nil), false, "no clear verdict, assumed real"},
		{"transientish word", "c", say("transiently odd\n", nil), false, "no clear verdict, assumed real"},
		{"empty output", "c", say("", nil), false, "no clear verdict, assumed real"},
		{"failing command", "c", say("transient: x", errors.New("exit 2")), false, "classifier failed: exit 2"},
		{"off", "  ", say("transient: x", nil), false, "classifier off"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v := ClassifyFailure(tc.command, "out", tc.shell)
			if v.Transient != tc.transient || v.Reason != tc.reason {
				t.Errorf("got %+v, want transient=%v reason=%q", v, tc.transient, tc.reason)
			}
		})
	}
}

func TestClassifyFailureGivesOutputAndCommand(t *testing.T) {
	var gotCmd, gotIn string
	ClassifyFailure("mycmd", "the output", func(c, in string) (string, error) {
		gotCmd, gotIn = c, in
		return "real: x", nil
	})
	if gotCmd != "mycmd" || gotIn != "the output" {
		t.Errorf("shell got %q, %q", gotCmd, gotIn)
	}
}
