package api

// This file holds the failure classifier: a machine-wide command that reads
// a failed step's output on standard input and prints whether the failure is
// transient (worth a retry) or real, with a one-line reason. The default
// runs Claude headless with Haiku; any LLM CLI can stand in, and an empty
// setting turns classification off.

import (
	"fmt"
	"regexp"
	"strings"
)

// FailureClassifierKey is the setting naming the classifier command.
const FailureClassifierKey = "rhiza.daemon.failureClassifier"

// DefaultFailureClassifier runs Claude headless with Haiku.
const DefaultFailureClassifier = `claude -p --model haiku "You classify a failed command from a software build or CI run. Its output follows on standard input. Answer with one line: 'transient: <reason>' if the failure is an infrastructure or environment hiccup that a plain rerun would likely fix (a network or GitHub error, a held lock, a vanished temp file, a full or flaky runner), or 'real: <reason>' if the code, a test, or the configuration is at fault. Give only that line."`

// FailureVerdict is the classifier's reading of one failure.
type FailureVerdict struct {
	Transient bool
	Reason    string
}

var verdictRE = regexp.MustCompile(`(?i)^\s*(transient|real)\b\s*[:\-–—]?\s*(.*)$`)

// ParseFailureVerdict reads a classifier's output: its first non-empty line
// must start with transient or real. Anything else, empty output included,
// is real: only a clear transient allows a retry.
func ParseFailureVerdict(out string) FailureVerdict {
	for _, l := range strings.Split(out, "\n") {
		if strings.TrimSpace(l) == "" {
			continue
		}
		m := verdictRE.FindStringSubmatch(l)
		if m == nil {
			break
		}
		return FailureVerdict{Transient: strings.EqualFold(m[1], "transient"), Reason: strings.TrimSpace(m[2])}
	}
	return FailureVerdict{Reason: "no clear verdict, assumed real"}
}

// ClassifyFailure runs command (through sh -c, by shell) with output on its
// standard input and parses what it prints. An empty command classifies
// nothing: the failure is real. A failing command is real too.
func ClassifyFailure(command, output string, shell func(command, stdin string) (string, error)) FailureVerdict {
	if strings.TrimSpace(command) == "" || shell == nil {
		return FailureVerdict{Reason: "classifier off"}
	}
	out, err := shell(command, output)
	if err != nil {
		return FailureVerdict{Reason: fmt.Sprintf("classifier failed: %v", err)}
	}
	return ParseFailureVerdict(out)
}

// FailureClassifierCommand is the classifier command in effect for ctx: the
// setting's value, "" when it is turned off.
func FailureClassifierCommand(ctx *Ctx) string {
	rows, err := settingsOf(ctx, nil)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(settingValue(rows, FailureClassifierKey))
}
