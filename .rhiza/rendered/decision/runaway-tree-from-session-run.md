# runaway-tree-from-session-run

decision · `01M7QN2R8384ZPCVR2121MZ1SZ`

## statement

A session's process tree is the 'rhiza session run <session-id>' process and every descendant, read from one ps -A -o pid=,ppid=,rss=,args= call (Linux and macOS alike, rss in KiB on both), never /proc; the reading keeps each process (pid, ppid, rss, args), not only the totals, so the stall pass can later tell an idle agent with a live background command from one with nothing running (orchestrator evidence 2026-11-18, ~/orchestrator/planner-notes/m40-stalls.md)

## wording

(assumed, would have escalated) planner's choice under the autonomous rules, from the code and milestone 40's Decision m40-host-limits-as-holds

## Links

### Out

- derived-from [runaway-plan-notes](../source/runaway-plan-notes.md)

### In

- decided-by [session-process-tree-reading](../capability/rhiza/orchestration/daemon/daemon-watchdog/session-process-tree-reading.md)

## History

- `01M7NVMDMBVBKRW5B8SB0KMW6T` The daemon's watchdog treats a runaway agent session of any workflow, one whose process tree's process count or resident memory grows past a host threshold read portably (ps, on Linux and macOS), like a stall: it nudges the session, then fails the step with the reading as its reason
