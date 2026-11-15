package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
)

// replayEvent is the part of a go test -json event the replay reads.
type replayEvent struct {
	Action  string
	Package string
	Test    string
	Elapsed float64
	Output  string
}

// ErrReplayFailed is returned by Replay when go test failed: a failed test, a
// failed build or package, or a package that never ended (a killed go test).
var ErrReplayFailed = errors.New("go test failed")

// Replay reads a go test -json stream and writes what a CI log needs of it:
// go test's text for every test that did not pass (its output, ending in its
// --- FAIL: line or a panic), in the order they finished; the package's own
// text (its FAIL or ok line) and any build output as they arrive; then one
// script-time line per passing TestScripts subtest, sorted by name. A passing
// or skipped test's text is dropped. Lines that are not JSON pass through.
//
// It returns ErrReplayFailed (after writing everything) when a test, a build
// or a package failed, or when the stream ends with no package end event or
// with a started package unfinished.
func Replay(r io.Reader, out io.Writer) error {
	w := bufio.NewWriter(out)
	defer w.Flush()
	type key struct{ pkg, test string }
	text := map[key]*strings.Builder{} // output of tests still running
	var order []key                    // tests in the order they started
	started := map[string]bool{}       // packages started and not ended
	ended := 0
	failed := false
	times := map[string]float64{}

	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for sc.Scan() {
		line := sc.Bytes()
		var ev replayEvent
		if trimmed := bytes.TrimSpace(line); len(trimmed) == 0 || trimmed[0] != '{' || json.Unmarshal(trimmed, &ev) != nil || ev.Action == "" {
			w.Write(line)
			w.WriteByte('\n')
			continue
		}
		switch ev.Action {
		case "build-output":
			w.WriteString(ev.Output)
			continue
		case "build-fail":
			failed = true
			continue
		}
		if ev.Test == "" {
			// The package's own events.
			switch ev.Action {
			case "start":
				started[ev.Package] = true
			case "output":
				w.WriteString(ev.Output)
			case "pass", "skip":
				delete(started, ev.Package)
				ended++
			case "fail":
				delete(started, ev.Package)
				ended++
				failed = true
			}
			continue
		}
		k := key{ev.Package, ev.Test}
		switch ev.Action {
		case "output":
			b := text[k]
			if b == nil {
				b = &strings.Builder{}
				text[k] = b
				order = append(order, k)
			}
			b.WriteString(ev.Output)
		case "pass":
			delete(text, k)
			if name, ok := strings.CutPrefix(ev.Test, "TestScripts/"); ok && validName(name) {
				times[name] = math.Max(times[name], ev.Elapsed)
			}
		case "skip":
			delete(text, k)
		case "fail":
			failed = true
			if b := text[k]; b != nil {
				w.WriteString(b.String())
			}
			delete(text, k)
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}

	// A killed go test leaves tests running: their text is all there is.
	for _, k := range order {
		if b := text[k]; b != nil {
			w.WriteString(b.String())
			delete(text, k)
			failed = true
		}
	}
	if len(started) > 0 || ended == 0 {
		pkgs := make([]string, 0, len(started))
		for p := range started {
			pkgs = append(pkgs, p)
		}
		sort.Strings(pkgs)
		if len(pkgs) == 0 {
			fmt.Fprintln(w, "scripttimings: no package end event in the go test -json input (go test killed?)")
		} else {
			fmt.Fprintf(w, "scripttimings: no package end event for %s (go test killed?)\n", strings.Join(pkgs, ", "))
		}
		failed = true
	}

	names := make([]string, 0, len(times))
	for n := range times {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		w.WriteString(FormatScriptTime(n, times[n]))
		w.WriteByte('\n')
	}
	if failed {
		return ErrReplayFailed
	}
	return nil
}
