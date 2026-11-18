# session-process-tree-reading-qa-1

evidence · `01M7QPPEJV4FNHDQDMWMZZAPKX` · actor [claude-code-9dc16082](../actor/claude-code-9dc16082.md)

## result

PASS

## statement

QA at b9129e145: go test -trimpath -p 2 -count=1 ./internal/proctree -v exit 0 (TestSessionTree, KeepsEachProcess, PrefixIsNotRoot, NotFound, ParseSkipsHeaderAndMalformed, ReadPSFails, ReadUsesTable PASS). Real 'ps -A -ww -o pid= -o ppid= -o rss= -o args=' on this Linux host prints live session roots as '<pid> <ppid> <rss> /home/rhiza/bin/rhiza session run <uuid>', the shape IsSessionRoot matches (fields 1-3 session/run/id, whole-id match). Code read: one ps call, no /proc, no OS-specific call, not-found is an error, rss KiB x 1024.

## Links

### Out

- verifies [session-process-tree-reading](../capability/rhiza/orchestration/daemon/daemon-watchdog/session-process-tree-reading.md)

## History

- `01M7NVMDMBVBKRW5B8SB0KMW6T` The daemon's watchdog treats a runaway agent session of any workflow, one whose process tree's process count or resident memory grows past a host threshold read portably (ps, on Linux and macOS), like a stall: it nudges the session, then fails the step with the reading as its reason
