# session-process-tree-reading-verified

evidence · `01M7QN8J4GXF0WF77QMJ6K04FJ` · actor [claude-code-bf9f5d32](../actor/claude-code-bf9f5d32.md)

## result

PASS

## statement

At 7b871dabb re-ran the Checks block: go test -trimpath -p 2 -count=1 ./internal/proctree exit 0 (7 tests PASS), go vet ./internal/proctree exit 0, rhiza ship generate exit 0, rhiza render exit 0 (only the verifier's own actor rendered). Read the diff: tests use fake ps text via Parse and a swapped runPS, never the host's ps; the table holds sessA (200..204, three deep below the root) and sessB (300..303) trees, each giving its own pids, count and rss summed as KiB x 1024 with the root counted and the other session's, prefix-sharing sessA2's, a Rex tab's and unrelated processes left out; 'sess' (a prefix) gives ErrNotFound and IsSessionRoot rejects sessA2 for sessA; an absent session and an empty table give ErrNotFound naming the session, not a zero tree; the header, garbage and bad-number lines are skipped (15 parsed); a failing ps gives an error starting 'ps: '. Only internal/proctree and the builder's actor record changed: no /proc read (a doc comment only), no OS-specific calls, no new dependency, internal/api and the watchdog untouched, gofmt clean, no TODOs.

## Links

### Out

- verifies [session-process-tree-reading](../capability/rhiza/orchestration/daemon/daemon-watchdog/session-process-tree-reading.md)

## History

- `01M7NVMDMBVBKRW5B8SB0KMW6T` The daemon's watchdog treats a runaway agent session of any workflow, one whose process tree's process count or resident memory grows past a host threshold read portably (ps, on Linux and macOS), like a stall: it nudges the session, then fails the step with the reading as its reason
