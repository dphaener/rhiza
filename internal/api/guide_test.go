package api

import (
	"bytes"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/genesis/internal/guide"
)

// pendingCommands are commands the guide names before they exist. Tasks 7
// and 8 remove the reads (slice, why, log, roadmap, ready, show, tree, ls)
// as they land.
var pendingCommands = []string{"why", "log", "roadmap", "ready", "show", "tree", "ls"}

// commandWord is one word of a command path.
var commandWord = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// codeSpans returns the code in markdown: each fenced block's body and each
// inline backtick span outside the fences.
func codeSpans(md string) []string {
	var spans, prose []string
	var fence []string
	inFence := false
	for line := range strings.SplitSeq(md, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			if inFence {
				spans = append(spans, strings.Join(fence, "\n"))
				fence = nil
			}
			inFence = !inFence
			continue
		}
		if inFence {
			fence = append(fence, line)
		} else {
			prose = append(prose, line)
		}
	}
	for i, part := range strings.Split(strings.Join(prose, "\n"), "`") {
		if i%2 == 1 {
			spans = append(spans, part)
		}
	}
	return spans
}

// commandRefs returns the command references in md's code: for each
// "genesis" token, the command words after it, up to the first token that
// is not one. A bare genesis is skipped.
func commandRefs(md string) [][]string {
	var refs [][]string
	for _, span := range codeSpans(md) {
		f := strings.Fields(span)
		for i, tok := range f {
			if tok != "genesis" {
				continue
			}
			var words []string
			for _, w := range f[i+1:] {
				if !commandWord.MatchString(w) {
					break
				}
				words = append(words, w)
			}
			if len(words) > 0 {
				refs = append(refs, words)
			}
		}
	}
	return refs
}

// resolves reports whether words name a command as the CLI resolves a
// command line: some prefix is an operation, or all of words name a group
// of operations (genesis change). A group followed by a word that names
// none of its operations (genesis change frobnicate) does not resolve.
func resolves(words []string) bool {
	for n := len(words); n > 0; n-- {
		name := strings.Join(words[:n], ".")
		if _, ok := Lookup(name); ok {
			return true
		}
		if n < len(words) {
			continue
		}
		for _, o := range Operations() {
			if strings.HasPrefix(o.Name, name+".") {
				return true
			}
		}
	}
	return false
}

// badRefs returns the references in md that resolve to no operation and are
// not pending, as command lines.
func badRefs(md string) []string {
	var bad []string
	for _, words := range commandRefs(md) {
		if resolves(words) || slices.Contains(pendingCommands, words[0]) {
			continue
		}
		bad = append(bad, "genesis "+strings.Join(words, " "))
	}
	return bad
}

func TestGuideCommandsExist(t *testing.T) {
	for _, txt := range guide.List() {
		if len(commandRefs(txt.Body)) == 0 {
			t.Errorf("guide %s: no command references found; is the scanner broken?", txt.Name)
		}
		for _, ref := range badRefs(txt.Body) {
			t.Errorf("guide %s: %q names no operation", txt.Name, ref)
		}
	}
}

func TestPendingCommandsUnregistered(t *testing.T) {
	for _, c := range pendingCommands {
		if resolves([]string{c}) {
			t.Errorf("pending command %q is registered now; remove it from pendingCommands", c)
		}
	}
}

func TestCommandRefsScanner(t *testing.T) {
	md := "Write only through `genesis` commands; see `genesis-scrub`.\n" +
		"Run `genesis frobnicate --now`, then `genesis change new \"x\"`,\n" +
		"`genesis link <node> derived-from <h>` and `genesis guide scrub`.\n" +
		"`genesis change frobnicate` is a typo; bare `genesis change` lists the group.\n" +
		"Plain genesis frobnicate outside code is prose.\n" +
		"```sh\n$ genesis describe apply --json\ngenesis slice <ref>\ngenesis zork it\n```\n"
	got := commandRefs(md)
	want := [][]string{
		{"describe", "apply"}, {"slice"}, {"zork", "it"},
		{"frobnicate"}, {"change", "new"}, {"link"}, {"guide", "scrub"},
		{"change", "frobnicate"}, {"change"},
	}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Errorf("commandRefs = %q, want %q", got, want)
	}
	bad := badRefs(md)
	if wantBad := []string{"genesis zork it", "genesis frobnicate", "genesis change frobnicate"}; !slices.Equal(bad, wantBad) {
		t.Errorf("badRefs = %q, want %q", bad, wantBad)
	}
}

func TestGuideNamesAreNotSubcommands(t *testing.T) {
	for _, o := range Operations() {
		sub, ok := strings.CutPrefix(o.Name, "guide.")
		if !ok {
			continue
		}
		sub, _, _ = strings.Cut(sub, ".")
		if _, ok := guide.Lookup(sub); ok {
			t.Errorf("guide text %q is shadowed by the operation %s", sub, o.Name)
		}
	}
}

func TestGuide(t *testing.T) {
	ctx := &Ctx{Version: "v9"}
	got, err := runGuide(ctx, GuideRequest{})
	if err != nil {
		t.Fatal(err)
	}
	want, _ := guide.Lookup("genesis")
	r := got.(GuideResponse)
	if r.Name != "genesis" || r.Version != "v9" || r.Text != want.Body {
		t.Errorf("guide = %+v", r)
	}
	var buf bytes.Buffer
	r.Human(&buf)
	if buf.String() != want.Body {
		t.Error("Human output differs from the text")
	}
	_, err = runGuide(ctx, GuideRequest{Name: "nope"})
	if e := AsError(err); e.Code != CodeNotFound || e.Next != "genesis guide ls" {
		t.Errorf("unknown name: %+v", e)
	}
}

func TestGuideLs(t *testing.T) {
	got, err := runGuideLs(&Ctx{}, GuideLsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	r := got.(GuideLsResponse)
	if len(r.Texts) != len(guide.List()) {
		t.Fatalf("guide ls: %d texts, want %d", len(r.Texts), len(guide.List()))
	}
	for i, txt := range guide.List() {
		info := r.Texts[i]
		kind := map[bool]string{true: "skill", false: "topic"}[txt.Skill]
		if info.Name != txt.Name || info.Kind != kind || info.Tokens != guide.Tokens(txt.Body) || info.Description != txt.Description {
			t.Errorf("guide ls [%d] = %+v", i, info)
		}
	}
}
