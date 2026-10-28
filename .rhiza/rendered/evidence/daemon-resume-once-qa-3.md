# daemon-resume-once-qa-3

evidence · `01M62D0FACXQA0ERW6NWVVAGND` · actor [claude-code-73eca481](../actor/claude-code-73eca481.md)

## result

PASS

## statement

QA cycle 3 with the built binary (359e3c1) in a scratch repo, with a fake rex that logs every call and runs only the daemon's own `rhiza run start|resume` tab command (agent tabs only logged): killed a daemon-started run with kill -9; the next `rhiza daemon --once` opened one tab running 'run resume <run>', journaled daemon-resumed and the run continued build/build; a poll while it ran opened no tab; after killing it again a poll opened no tab, daemon-resumed stayed at 1, and `rhiza run start` refused the change set naming the live run.

## Links

### Out

- verifies [daemon-resume-once](../capability/rhiza/orchestration/daemon/daemon-resume-once.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
