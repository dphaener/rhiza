# daemon-status-qa-3

evidence · `01M62D0FRSDPF1ZR9A7K2WHPHK` · actor [claude-code-73eca481](../actor/claude-code-73eca481.md)

## result

PASS

## statement

QA cycle 3 with the built binary (359e3c1) in a scratch repo, with a fake rex that logs every call and runs only the daemon's own `rhiza run start|resume` tab command (agent tabs only logged): `rhiza daemon status` said 'Daemon not running' before any daemon; while one ran it named its process; after runs it listed Live runs (interrupted / running with change sets), Last started, Last resumed and Not journaled, and --json carried the same fields; after kill -9 of the daemon it said not running.

## Links

### Out

- verifies [daemon-status](../capability/rhiza/orchestration/daemon/daemon-status.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
