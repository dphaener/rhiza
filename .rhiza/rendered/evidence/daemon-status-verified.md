# daemon-status-verified

evidence · `01M629NQ39RM039MAGQYTCV51F` · actor [claude-code-14c860c0](../actor/claude-code-14c860c0.md)

## result

PASS

## statement

At af6b1ab, go test ./cmd/rhiza -run 'TestScripts/(daemon-status|describe-golden)$' -count=1 -v passes. daemon-status.txtar runs rhiza daemon status before any daemon (text 'Daemon not running', --json running:false); after rhiza daemon --once --max-runs=3 --autonomous with a fake rex on PATH it shows the state the daemon left: text 'Daemon not running' plus last daemon process/stopped, Max runs: 3, Autonomous: true, Last started: a run of <change>, and --json with running:false, max_runs 3, autonomous true, last_started {change, at}; then with the state file removed and daemon.lock naming pid 999999 it says not running and no Last started. describe-golden.txtar carries the new daemon.status operation and passes. Diff read: daemon.State written under the git common dir (.git/rhiza/daemon.state.json) at each poll and on stop, with process, start, last poll, max runs, autonomous, live runs and last started/resumed; status takes running from daemon.Running on the lock. make fmt-check vet test pass; make self-check fails only on .rhiza/rendered actor/index drift that HEAD~1 shows too, outside this leaf.

## Links

### Out

- verifies [daemon-status](../capability/rhiza/orchestration/daemon/daemon-status.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
