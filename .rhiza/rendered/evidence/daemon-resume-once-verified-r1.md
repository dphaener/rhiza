# daemon-resume-once-verified-r1

evidence · `01M62CECJMW7AHK0E6GTM3WZFR` · actor [claude-code-89e205f8](../actor/claude-code-89e205f8.md)

## result

PASS

## statement

At cc732b0 (review fix of quality f1/f2): read c5e4d67 and cc732b0. Resumable now skips a run whose store lock another holder has (the daemon passes rs.Live, a non-blocking flock probe) and a run with no process record started under Unrecorded (1m) ago, covering the gap between journaling a run and its process entry, before and after the lock is taken. go test -count=1 ./internal/daemon -run Resumable -v: TestResumable (interrupted, resumed-again, failed, starting, unrecorded, locked) and TestResumableLeavesAStoreLockedRun (real Store.TryLock held: no resume; released: R1 resumed) PASS. go test -count=1 ./cmd/rhiza -run TestScripts/daemon PASS: daemon-resume.txtar runs rhiza daemon --once twice with a fake Rex over a fixture journal naming dead pid 999999; the first opens one 'rex window new ... run resume INTERRUPTED1' tab and journals one daemon-resumed, the second resumes none, and FAILED00001 never appears in rex.log. make fmt-check vet test exit 0. make self-check fails only on stale .rhiza/rendered actor/index files from the run's own evidence commits, not from this leaf's code.

## Links

### Out

- verifies [daemon-resume-once](../capability/rhiza/orchestration/daemon/daemon-resume-once.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
