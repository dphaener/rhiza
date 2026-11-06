# daemon-lock-qa-3

evidence · `01M62D0DYBFS1PH4F35EX42QKF` · actor [claude-code-73eca481](../actor/claude-code-73eca481.md)

## result

PASS

## statement

QA cycle 3 with the built binary (359e3c1) in a scratch repo, with a fake rex that logs every call and runs only the daemon's own `rhiza run start|resume` tab command (agent tabs only logged): a daemon held the lock; `rhiza daemon --once` from the repo's run worktree (repo.built) exited 2 with 'another daemon holds the lock: process Mac.localdomain:88034'; after kill -9 of a running daemon, status said not running and a new `rhiza daemon --once` took the lock and ran.

## Links

### Out

- verifies daemon-lock (retired)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
