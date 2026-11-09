package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rogpeppe/go-internal/testscript"
)

// quickstartDoc is the script command `quickstart-doc [file]`: it runs every
// command in the ```bash blocks of docs/quickstart.md (or file), in order,
// in the script's work dir, and fails on the first one that fails. A line
// is `rhiza ...` or `git ...` (words may be single- or double-quoted), or a
// `cat > file <<'EOF'` heredoc that writes a file. Anything else in a bash
// block fails the script, so the doc cannot drift into commands the script
// does not run. Prose and sample output go in ```text blocks.
func quickstartDoc(ts *testscript.TestScript, neg bool, args []string) {
	file := "../../docs/quickstart.md"
	if neg || len(args) > 1 {
		ts.Fatalf("usage: quickstart-doc [file]")
	}
	if len(args) == 1 {
		file = args[0]
	}
	data, err := os.ReadFile(file)
	if err != nil {
		ts.Fatalf("quickstart-doc: %v", err)
	}
	var lines []string
	inBash, other := false, false
	for _, l := range strings.Split(string(data), "\n") {
		switch {
		case strings.HasPrefix(l, "```"):
			if inBash || other {
				inBash, other = false, false
			} else if strings.TrimPrefix(l, "```") == "bash" {
				inBash = true
			} else {
				other = true
			}
		case inBash:
			lines = append(lines, l)
		}
	}
	ran := 0
	for i := 0; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		if strings.HasPrefix(l, "cat > ") && strings.HasSuffix(l, " <<'EOF'") {
			name := strings.TrimSuffix(strings.TrimPrefix(l, "cat > "), " <<'EOF'")
			var body []string
			i++
			for ; i < len(lines) && lines[i] != "EOF"; i++ {
				body = append(body, lines[i])
			}
			if i == len(lines) {
				ts.Fatalf("quickstart-doc: heredoc for %s has no EOF", name)
			}
			if err := os.WriteFile(ts.MkAbs(name), []byte(strings.Join(body, "\n")+"\n"), 0o644); err != nil {
				ts.Fatalf("quickstart-doc: %v", err)
			}
			ran++
			continue
		}
		words, err := splitWords(l)
		if err != nil || len(words) == 0 || (words[0] != "rhiza" && words[0] != "git") {
			ts.Fatalf("quickstart-doc: cannot run %q (only rhiza, git and cat heredocs): %v", l, err)
		}
		if err := ts.Exec(words[0], words[1:]...); err != nil {
			ts.Fatalf("quickstart-doc: %s: %v", l, err)
		}
		ran++
	}
	if ran == 0 {
		ts.Fatalf("quickstart-doc: no commands in %s", file)
	}
	fmt.Fprintf(ts.Stdout(), "quickstart-doc: ran %d commands\n", ran)
}

// splitWords splits a command line into words, honouring single and double
// quotes. It does no expansion.
func splitWords(s string) ([]string, error) {
	var words []string
	var cur strings.Builder
	inWord := false
	var quote rune
	for _, r := range s {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				cur.WriteRune(r)
			}
		case r == '\'' || r == '"':
			quote, inWord = r, true
		case r == ' ' || r == '\t':
			if inWord {
				words = append(words, cur.String())
				cur.Reset()
				inWord = false
			}
		default:
			cur.WriteRune(r)
			inWord = true
		}
	}
	if quote != 0 {
		return nil, fmt.Errorf("unterminated quote")
	}
	if inWord {
		words = append(words, cur.String())
	}
	return words, nil
}
