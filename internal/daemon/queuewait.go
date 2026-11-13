package daemon

import (
	"bufio"
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// QueueWaitFile is the file, under Dir, that the daemon appends why an
// accepted change set waits before its run starts to, one JSON line each.
// The perf analysis reads it directly.
const QueueWaitFile = "queue.wait.jsonl"

// QueueWaitStarted is the Reason of the line written when the daemon starts
// a change set's run.
const QueueWaitStarted = "started"

// QueueWaitOutside is the Detail of a started line for a run the daemon did
// not start: one started by hand, or an owner-shipped change set's.
const QueueWaitOutside = "started outside the daemon"

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
// Title set). First, for each change set whose last journaled line is a
// pre-start reason, a run in runs that started after that line (a run the
// daemon did not start: by hand, or a ship-owner change set's) is journaled
// as started at the run's start time with QueueWaitOutside as its detail,
// whatever the poll's filters say. A skip is journaled only for a pre-start
// reason, only when states has no post-start state for its change set (so a
// hold that comes on while its run is live or held writes nothing), and only
// when it differs from the change set's last journaled reason, which is read
// from the journal, so a restarted daemon writes no duplicate. A start is
// always journaled. It returns the first error and goes on writing.
func JournalQueueWaits(common string, at time.Time, states *RunStates, runs []run.Journal, skips, started []Skip) error {
	last := map[string]QueueWait{}
	for _, w := range ReadQueueWaits(common) {
		last[w.Change] = w
	}
	var first error
	record := func(w QueueWait) {
		if err := RecordQueueWait(common, w); err != nil && first == nil {
			first = err
		}
		last[w.Change] = w
	}
	outside := map[string]time.Time{}
	for _, j := range runs {
		st := j.State()
		w, ok := last[st.Change]
		if st.Change == "" || !ok || !preStartReasons[w.Reason] || !st.StartedAt.After(w.At) {
			continue
		}
		if t, ok := outside[st.Change]; !ok || st.StartedAt.Before(t) {
			outside[st.Change] = st.StartedAt
		}
	}
	for _, cs := range slices.Sorted(maps.Keys(outside)) {
		record(QueueWait{At: outside[cs].UTC(), Change: cs, Title: last[cs].Title, Reason: QueueWaitStarted, Detail: QueueWaitOutside})
	}
	write := func(sk Skip, reason string) {
		record(QueueWait{At: at.UTC(), Change: sk.Change, Title: sk.Title, Reason: reason, Detail: strings.Join(sk.Detail, "; ")})
	}
	for _, sk := range skips {
		if sk.Change == "" || !preStartReasons[sk.Reason] || last[sk.Change].Reason == sk.Reason || states.PostStart(sk.Change) != "" {
			continue
		}
		write(sk, sk.Reason)
	}
	for _, sk := range started {
		write(sk, QueueWaitStarted)
	}
	return first
}
