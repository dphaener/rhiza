# maxruns-auto-poll-verified

evidence · `01M7QYCZTEXQBJ5FHC9A46R90Y` · actor [claude-code-f9f94796](../actor/claude-code-f9f94796.md)

## result

PASS

## statement

At d1df56a6c re-ran every Checks: command, all exit 0 (capped scopes, -p 2): go test ./cmd/rhiza -run TestScripts/daemon-maxruns-auto (fresh, -count=1) and the ten daemon/describe scripts; go test ./internal/api -run MaxRuns (TestDaemonMaxRunsAutoCap, TestDaemonMaxRunsAuto pass); go test ./internal/api/... and ./internal/daemon/... ok; rhiza ship generate changed 0 files; rhiza render left the committed rendered files unchanged. daemon-maxruns-auto.txtar covers the acceptance: with auto and idle readings the first --once poll starts exactly one of three accepted change sets and the next poll starts a second; daemon.state.json has max_runs_auto true with max_runs 1 then 2; maxRuns=1 starts none; RHIZA_TEST_LOAD1/5=3.5 starts none, and so does the tmp-low hold; with both cleared the third run starts. The diff works out the cap as daemonConf.autoMaxRuns(running): running+1 with room, running while the machine hold, tmp-low hold or disk-low check holds. It uses the poll's own readings and reads nothing anew. A numeric maxRuns is unchanged. The edits outside the leaf's files list are needed: daemon_projects.go moves the existing disk-low read into prepare so the cap and demand share one reading (required by must_not), and the two test files only update poll() calls for the new parameter. No stubs or TODOs.

## Links

### Out

- verifies [maxruns-auto-poll](../capability/rhiza/orchestration/daemon/daemon-picks/maxruns-auto-poll.md)

## History

- `01M7NVMXZ1YEGZB56ZEAMPX44B` rhiza.daemon.maxRuns=auto lets the daemon start another run only while the pre-start checks (load, memory, free tmp, disk) have room; a number still pins the limit
