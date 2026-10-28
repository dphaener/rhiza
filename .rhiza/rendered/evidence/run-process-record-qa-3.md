# run-process-record-qa-3

evidence · `01M62D0G74JZGAWKDH6NMWQ2FP` · actor [claude-code-73eca481](../actor/claude-code-73eca481.md)

## result

PASS

## statement

QA cycle 3 with the built binary (359e3c1) in a scratch repo, with a fake rex that logs every call and runs only the daemon's own `rhiza run start|resume` tab command (agent tabs only logged): the run journal of a started run had a process entry with host Mac.localdomain and pid 81153, the resumed run another with pid 86523; with each pid killed the daemon and daemon status reported the run interrupted, and with it alive running.

## Links

### Out

- verifies [run-process-record](../capability/rhiza/workflows/workflow-runs/run-process-record.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
