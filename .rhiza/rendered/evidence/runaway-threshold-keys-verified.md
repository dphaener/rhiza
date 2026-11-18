# runaway-threshold-keys-verified

evidence · `01M7QNGMW1NJ7CYE6EN7HRMFGH` · actor [claude-code-23a30f3d](../actor/claude-code-23a30f3d.md)

## result

PASS

## statement

Verified commit 6e7b448d0. go test -trimpath -p 2 ./internal/api -run TestWatchdogPolicy: exit 0 (runaway keys subtest resolves defaults 400, 8<<30 bytes, 5m and project over machine; invalid-values subtest refuses runawayMemory=lots, runawayProcs=-1, runawayGrace=0 naming the key). env -u RHIZA_*/SHIP_*/REX_* go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/watchdog-policy$': exit 0 (three keys listed at defaults for machine and project, runawayMemory 12G set via rhiza daemon set, 'lots' and '-- runawayProcs -1' refused naming the key). rhiza ship generate: exit 0, 0 changed files (docs/reference/config.md current). rhiza render: exit 0. make self-check: exit 0 from a clean detached worktree at the commit (in the run worktree it only failed on this verifier's own untracked actor entry). Diff touches only the leaf's files plus its .rhiza log/rendered; existing watchdog keys unchanged; no process table read; no TODOs.

## Links

### Out

- verifies [runaway-threshold-keys](../capability/rhiza/orchestration/daemon/daemon-watchdog/runaway-threshold-keys.md)

## History

- `01M7NVMDMBVBKRW5B8SB0KMW6T` The daemon's watchdog treats a runaway agent session of any workflow, one whose process tree's process count or resident memory grows past a host threshold read portably (ps, on Linux and macOS), like a stall: it nudges the session, then fails the step with the reading as its reason
