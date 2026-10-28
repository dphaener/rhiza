# daemon-status-qa-1

evidence · `01M629Y7QNZHSQW79FF21V2CX9` · actor [claude-code-715b2246](../actor/claude-code-715b2246.md)

## result

PASS

## statement

Built rhiza at f087c0e into $TMPDIR/rhiza-qa and drove it by hand in a scratch repo (RHIZA_* env unset, rex shimmed on PATH): rhiza daemon status before any daemon printed 'Daemon not running' (--json running false, runs []); while a daemon ran it printed 'Daemon running, process Mac.localdomain:89202' with Started, Last poll, Max runs 2, Autonomous false and Last started (also in --json); after kill -INT it printed not running plus 'Last daemon: process ..., stopped ...'; after the journal polls it listed the live runs as running or interrupted with their change sets and Last resumed; a lock naming dead pid 999999 read as not running.

## Links

### Out

- verifies [daemon-status](../capability/rhiza/orchestration/daemon/daemon-status.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
