# daemon

capability · `01M61R7HFGRYD2Q1FG4K5N65ME` · parent [orchestration](../orchestration.md)

## statement

rhiza daemon orchestrates a repo's milestone work: it starts ship-change runs on accepted, ready change sets under a cap, resumes interrupted runs once and reports its state, with no socket and no LLM conductor.

## Links

### Out

- decided-by [answer-8cd68bc7-2](../../../decision/answer-8cd68bc7-2.md)
- decided-by [answer-8cd68bc7-3](../../../decision/answer-8cd68bc7-3.md)
- decided-by [answer-8cd68bc7-14](../../../decision/answer-8cd68bc7-14.md)

### In

- verifies [pr-52-the-daemon](../../../evidence/pr-52-the-daemon.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
