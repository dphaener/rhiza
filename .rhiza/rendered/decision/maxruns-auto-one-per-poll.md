# maxruns-auto-one-per-poll

decision · `01M7QXK2V41VD09860J5HJDD60`

## statement

With rhiza.daemon.maxRuns=auto each poll's machine-wide cap is the runs running plus one when every pre-start check (load and memory machine hold, tmp-low hold, the home project's disk-low check) has room, and the runs running when any holds, so each further start waits for the next poll's fresh readings

## wording

(planner ruling, derived from the intent and the daemon's pre-start checks; would have escalated)

## Links

### Out

- derived-from [maxruns-auto-plan](../source/maxruns-auto-plan.md)

### In

- decided-by [maxruns-auto-poll](../capability/rhiza/orchestration/daemon/daemon-picks/maxruns-auto-poll.md)
- decided-by [maxruns-auto-shown](../capability/rhiza/orchestration/daemon/daemon-status/maxruns-auto-shown.md)

## History

- `01M7NVMXZ1YEGZB56ZEAMPX44B` rhiza.daemon.maxRuns=auto lets the daemon start another run only while the pre-start checks (load, memory, free tmp, disk) have room; a number still pins the limit
