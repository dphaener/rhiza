# maxruns-auto-scope

decision · `01M7QXK2X24GT08P4VP8E709NH`

## statement

auto is a value of rhiza.daemon.maxRuns only: the default stays 2, a number still pins the cap, and rhiza.project.maxRuns and the daemon's --max-runs flag stay whole numbers

## wording

(planner ruling, derived from the intent's 'a number still pins the limit'; would have escalated)

## Links

### Out

- derived-from [maxruns-auto-plan](../source/maxruns-auto-plan.md)

### In

- decided-by [maxruns-auto-setting](../capability/rhiza/orchestration/daemon/daemon-settings/maxruns-auto-setting.md)

## History

- `01M7NVMXZ1YEGZB56ZEAMPX44B` rhiza.daemon.maxRuns=auto lets the daemon start another run only while the pre-start checks (load, memory, free tmp, disk) have room; a number still pins the limit
