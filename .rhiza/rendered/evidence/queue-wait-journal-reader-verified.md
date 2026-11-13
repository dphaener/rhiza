# queue-wait-journal-reader-verified

evidence · `01M79VCFRRG477YCWNYBHVAHHV` · actor [claude-code-2efc57f0](../actor/claude-code-2efc57f0.md)

## result

PASS

## statement

At eb1b5eb22 re-ran the Checks block: go test -trimpath ./internal/daemon/... -run QueueWait (exit 0, also -race -count=3 exit 0), rhiza ship generate (exit 0), rhiza render (exit 0, only this verifier's own actor rendered). Read the diff: internal/daemon/queuewait.go defines QueueWait (at, change, title, reason, detail), RecordQueueWait appends one marshalled line with a single O_APPEND write to queue.wait.jsonl under Dir(common), ReadQueueWaits reads lines oldest first skipping unparseable ones and returns nil on a missing file. queuewait_test.go covers each acceptance point: round trip with every field in order, a malformed line between two good ones skipped, missing file reads as no lines, 50 concurrent goroutine appends all read back whole and distinct. Nothing outside the leaf changed besides the builder's actor record; no stubs or TODOs; gofmt and go vet clean.

## Links

### Out

- verifies [queue-wait-journal-reader](../capability/rhiza/orchestration/daemon/daemon-picks/queue-wait-journal-reader.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
