# daemon-lock

capability · `01M61R7J7KWQ20F4X22KTJ924R` · parent [daemon](../daemon.md)

## acceptance

A testscript runs a daemon in the background with a fake Rex and starts a second from another worktree, which exits 2 naming the first; a lock file naming a dead process is taken over by a new daemon; no test in the default run binds a socket.

## statement

rhiza daemon runs in the foreground, one per repo: it holds a lock under the git common dir, so a second daemon started from any worktree of the repo exits 2 naming the first one's process, a lock left by a process that is gone is taken over, and the daemon binds no socket.

## Links

### Out

- decided-by [answer-8cd68bc7-14](../../../../decision/answer-8cd68bc7-14.md)

### In

- depends-on [daemon-picks](daemon-picks.md)
- depends-on [daemon-status](daemon-status.md)
- verifies [daemon-lock-verified](../../../../evidence/daemon-lock-verified.md)
- verifies [daemon-lock-qa-1](../../../../evidence/daemon-lock-qa-1.md)
- verifies [daemon-lock-qa-3](../../../../evidence/daemon-lock-qa-3.md)
- verifies [pr-52-the-daemon](../../../../evidence/pr-52-the-daemon.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
