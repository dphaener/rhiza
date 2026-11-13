# queue-wait-journal-reader-qa-2

evidence · `01M79XSHJQ4DQBQ5FWXDWSXR57` · actor [claude-code-72ff73af](../actor/claude-code-72ff73af.md)

## result

PASS

## statement

QA cycle 2: the built binary's daemon polls (separate processes) appended to <common>/rhiza/queue.wait.jsonl one whole JSON line per write with at, change, title, reason and detail, 12 lines over the session, all parseable. go test -trimpath -count=1 ./internal/daemon/... ok.

## Links

### Out

- verifies [queue-wait-journal-reader](../capability/rhiza/orchestration/daemon/daemon-picks/queue-wait-journal-reader.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
