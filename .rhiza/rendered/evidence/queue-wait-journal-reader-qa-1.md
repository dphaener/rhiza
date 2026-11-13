# queue-wait-journal-reader-qa-1

evidence · `01M79WHSY4308FW8K87VN02QQA` · actor [claude-code-037e23c0](../actor/claude-code-037e23c0.md)

## result

PASS

## statement

QA cycle 1: built this branch (go build ./cmd/rhiza) and ran it in a hermetic scratch Rhiza root (env -i, fake rex, RHIZA_TEST_DETACH_LOG). Real rhiza daemon --once polls wrote <common>/rhiza/queue.wait.jsonl, one JSON line per append with at, change, title, reason and detail; the file grew across separate daemon processes and every line parsed. go test -trimpath ./internal/daemon/... -run QueueWait ok.

## Links

### Out

- verifies [queue-wait-journal-reader](../capability/rhiza/orchestration/daemon/daemon-picks/queue-wait-journal-reader.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
