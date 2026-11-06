# daemon-lock-verified

evidence · `01M6279XQTZV46ZAYY08KFJXXP` · actor [claude-code-897dc5d6](../actor/claude-code-897dc5d6.md)

## result

PASS

## statement

At 2ad530b: go test -count=1 ./cmd/rhiza -run 'TestScripts/daemon-lock$' -v passed (and 8 more repeats with -count=8). The script runs rhiza daemon in the background with a fake rex on PATH; a second rhiza daemon from the repo and from a git worktree add'ed other worktree exits 2 with 'another daemon holds the lock: process <host>:<pid>' whose pid matches the first's lock file; after kill -INT the lock file is emptied and no daemon.sock exists; a lock file naming dead pid 999999 is taken over by a new daemon; rex.log is never written. internal/daemon unit tests pass under -race. grep finds no net.Listen/ListenUnix anywhere in the Go code, so no test binds a socket. Diff read: flock (LOCK_EX|LOCK_NB) on <git common dir>/rhiza/daemon.lock, holder JSON written after locking. make check exits 0.

## Links

### Out

- verifies daemon-lock (retired)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
