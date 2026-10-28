# daemon-status-needs-picks

decision · `01M620K0S92BPZKN58JJMVBCYS`

## statement

daemon-status depends on daemon-picks and daemon-run-tabs, so its check reads the state of a daemon that has picked and started a run. The check stays as written.

## wording

the orchestrator's option, picked by the owner: "Add the needs (Recommended)" (daemon-status needs daemon-picks and daemon-run-tabs; inbox-notify needs inbox-list; the checks stay as strong as written)

## Links

### Out

- answers [review-01m61h5z-r3-f1](../question/review-01m61h5z-r3-f1.md)

### In

- decided-by [daemon-status](../capability/rhiza/orchestration/daemon/daemon-status.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
