# daemon-lock-qa-1

evidence · `01M629Y5S11SRYKQ15PGF38FWT` · actor [claude-code-715b2246](../actor/claude-code-715b2246.md)

## result

PASS

## statement

Built rhiza at f087c0e into $TMPDIR/rhiza-qa and drove it by hand in a scratch repo (RHIZA_* env unset, rex shimmed on PATH): a long-running rhiza daemon held .git/rhiza/daemon.lock; a second rhiza daemon from a git worktree of the repo exited 2 with 'another daemon holds the lock: process Mac.localdomain:89202' (the first's pid) and next rhiza daemon status; lsof showed no unix or inet socket on the daemon; kill -INT stopped it with exit 0; a daemon.lock hand-written with dead pid 999999 read as not running and a new daemon took it over (exit 0).

## Links

### Out

- verifies daemon-lock (retired)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
