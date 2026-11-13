package daemon

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// QueueWaitFile is the file, under Dir, that the daemon appends why an
// accepted change set waits before its run starts to, one JSON line each.
// The perf analysis reads it directly.
const QueueWaitFile = "queue.wait.jsonl"

// QueueWaitStarted is the Reason of the line written when the daemon starts
// a change set's run.
const QueueWaitStarted = "started"

// QueueWait is one line of the queue-wait journal.
type QueueWait struct {
	// At is when the daemon wrote the line.
	At time.Time `json:"at"`
	// Change is the change set's ID and Title its title.
	Change string `json:"change"`
	Title  string `json:"title"`
	// Reason is the poll's pre-start Skip reason, or QueueWaitStarted.
	Reason string `json:"reason"`
	// Detail is the reason's detail text, empty when it has none.
	Detail string `json:"detail"`
}

// RecordQueueWait appends a line to the journal under a repo's git common
// dir.
func RecordQueueWait(common string, w QueueWait) error { return RecordQueueWaitIn(Dir(common), w) }

// ReadQueueWaits returns every journal line under a repo's git common dir.
func ReadQueueWaits(common string) []QueueWait { return ReadQueueWaitsIn(Dir(common)) }

// RecordQueueWaitIn appends a line to the journal in dir with a single
// write, so appends from several goroutines or processes do not interleave.
func RecordQueueWaitIn(dir string, w QueueWait) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	b, err := json.Marshal(w)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(dir, QueueWaitFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	_, err = f.Write(append(b, '\n'))
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	return err
}

// ReadQueueWaitsIn returns every line of the journal in dir, oldest first as
// written. Lines that do not parse are skipped, and a missing file reads as
// no lines.
func ReadQueueWaitsIn(dir string) []QueueWait {
	f, err := os.Open(filepath.Join(dir, QueueWaitFile))
	if err != nil {
		return nil
	}
	defer f.Close()
	var out []QueueWait
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadBytes('\n')
		var w QueueWait
		if len(line) > 0 && json.Unmarshal(line, &w) == nil && w.Change != "" {
			out = append(out, w)
		}
		if err != nil {
			return out
		}
	}
}

// preStartReasons are the Skip reasons of a change set that waits before its
// run starts: the only ones the journal records.
var preStartReasons = map[string]bool{
	SkipNoSlot: true, SkipNotReady: true, SkipNoMilestone: true, SkipNoLeaves: true,
	SkipDiskLow: true, SkipLoadHigh: true, SkipMemoryLow: true, SkipUsageReserve: true,
	SkipPaused: true, SkipShipOwner: true, SkipWorkflow: true,
}

// JournalQueueWaits appends the lines one poll owes the journal under a
// repo's git common dir. skips are the poll's skips, each with its Title
// filled in, and started the change sets whose runs the poll started (their
// Title set). A skip is journaled only for a pre-start reason, only when
// states has no post-start state for its change set (so a hold that comes on
// while its run is live or held writes nothing), and only when it differs
// from the change set's last journaled reason, which is read from the
// journal, so a restarted daemon writes no duplicate. A start is always
// journaled. It returns the first error and goes on writing.
func JournalQueueWaits(common string, at time.Time, states *RunStates, skips, started []Skip) error {
	last := map[string]string{}
	for _, w := range ReadQueueWaits(common) {
		last[w.Change] = w.Reason
	}
	var first error
	write := func(sk Skip, reason string) {
		w := QueueWait{At: at.UTC(), Change: sk.Change, Title: sk.Title, Reason: reason, Detail: strings.Join(sk.Detail, "; ")}
		if err := RecordQueueWait(common, w); err != nil && first == nil {
			first = err
		}
		last[sk.Change] = reason
	}
	for _, sk := range skips {
		if sk.Change == "" || !preStartReasons[sk.Reason] || last[sk.Change] == sk.Reason || states.PostStart(sk.Change) != "" {
			continue
		}
		write(sk, sk.Reason)
	}
	for _, sk := range started {
		write(sk, QueueWaitStarted)
	}
	return first
}
