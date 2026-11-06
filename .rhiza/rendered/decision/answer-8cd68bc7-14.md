# answer-8cd68bc7-14

decision · `01M61R3AM3GF3XCNZ5PT584EZ7`

## options

```json
[
  "As described (Recommended)",
  "Runs as child processes"
]
```

## picked

As described (Recommended)

## question

Topic 6b, how does the daemon run and launch work? I recommend: `rhiza daemon` runs in the foreground of whatever Rex tab you start it in, one per repo (a lock under the git common dir), polling the roadmap and run journals every few seconds with no socket. It starts each run as `rhiza run start <cs>` in a new Rex tab, so a run outlives a daemon restart. `--max-runs` defaults to 2, `--autonomous` passes through to every run it starts, and `rhiza daemon status` reads a state file.

## session

8cd68bc7-9786-4df4-8723-01a1831302a8

## statement

Topic 6b, how does the daemon run and launch work? I recommend: `rhiza daemon` runs in the foreground of whatever Rex tab you start it in, one per repo (a lock under the git common dir), polling the roadmap and run journals every few seconds with no socket. It starts each run as `rhiza run start <cs>` in a new Rex tab, so a run outlives a daemon restart. `--max-runs` defaults to 2, `--autonomous` passes through to every run it starts, and `rhiza daemon status` reads a state file. — As described (Recommended)

## wording

the agent's option, picked by the owner in session 8cd68bc7-9786-4df4-8723-01a1831302a8

## Links

### Out

- derived-from [milestone-4-conversation](../source/milestone-4-conversation.md)

### In

- decided-by [daemon](../capability/rhiza/orchestration/daemon.md)
- decided-by daemon-lock (retired)
- decided-by [run-process-record](../capability/rhiza/workflows/workflow-runs/run-process-record.md)
- decided-by daemon-run-tabs (retired)
- decided-by [daemon-status](../capability/rhiza/orchestration/daemon/daemon-status.md)

## History

- `01M61H5VY8MD7E7SVS0HQQHQJ0` Plan milestone 4 with the owner: its intent, its done points and the change sets that serve it, planned by Rhiza's plan workflow.
