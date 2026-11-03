# daemon-run-tabs-verified-3

evidence · `01M62CANGA4RXXDCV27TCSZSF4` · actor [claude-code-e8f87e90](../actor/claude-code-e8f87e90.md)

## result

PASS

## statement

At 608fcae, go test ./cmd/rhiza -run TestScripts/daemon -count=1 passed, and go test ./..., fmt-check and vet passed. In the kept work dir of daemon-picks.txtar, the fake Rex's rex.log shows rhiza daemon --permission-mode manual --once opened exactly one new tab (rex window new ... -c <repo> -- rhiza run start 01M5PWGQ8GB4618G4RH2TV0ZAH --permission-mode manual). That ID is the milestone change set Built's, and the line has no --autonomous; the leafless and backlog change sets got no tab. Under --once --autonomous the command ends with --permission-mode manual --autonomous. The diff of 161f5f5 (internal/api/daemon.go daemonStart) appends --autonomous only when the daemon has it, and a daemon with no mode or a bypass mode is refused before any tab opens. make self-check fails only on stale .rhiza/rendered actor/index files from the run's own model writes, not on this leaf's code.

## Links

### Out

- verifies daemon-run-tabs (retired)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
