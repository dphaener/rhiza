# daemon-status

capability · `01M61R7NEEK19RP2QEHFRQ6W8S` · parent [daemon](../daemon.md)

## acceptance

A testscript reads rhiza daemon status before any daemon (not running), after rhiza daemon --once with a fake Rex (the state with the run it started, text and --json), and with a lock naming a dead process (not running); describe-golden is regenerated.

## statement

The daemon keeps a state file under the git common dir (its process, start time, last poll, max runs, autonomous or not, the live runs and what it last started or resumed), and rhiza daemon status reads it, in text and --json, saying not running when no live daemon holds the lock.

## Links

### Out

- depends-on daemon-lock (retired)
- depends-on [daemon-picks](daemon-picks.md)
- depends-on daemon-run-tabs (retired)
- decided-by [answer-8cd68bc7-14](../../../../decision/answer-8cd68bc7-14.md)
- decided-by [daemon-status-needs-picks](../../../../decision/daemon-status-needs-picks.md)

### In

- blocks [review-01m61h5z-r3-f1](../../../../question/review-01m61h5z-r3-f1.md)
- verifies [daemon-status-verified](../../../../evidence/daemon-status-verified.md)
- verifies [daemon-status-qa-1](../../../../evidence/daemon-status-qa-1.md)
- verifies [daemon-status-qa-3](../../../../evidence/daemon-status-qa-3.md)
- verifies [pr-52-the-daemon](../../../../evidence/pr-52-the-daemon.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
