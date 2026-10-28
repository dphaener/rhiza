# journal-run-tab

decision · `01M620K1S3GGZB4B0H711KJVT9`

## statement

A run started in Rex journals its Rex window and block with its process, when it starts and each time it is resumed, so the inbox can focus a run's tab while it is open and say it is gone once it is closed.

## wording

the orchestrator's option, picked by the owner: "Journal the run's tab (Recommended)" (when a run starts in Rex, its journal records its window and block, alongside run-process-record's host and pid, so inbox-open focuses that tab when it's still open, and otherwise says it's gone)

## Links

### In

- decided-by [run-process-record](../capability/rhiza/workflows/workflow-runs/run-process-record.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
