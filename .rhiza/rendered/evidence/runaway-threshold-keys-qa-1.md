# runaway-threshold-keys-qa-1

evidence · `01M7QPPEQ6BWMBP0KFZHXD5FB3` · actor [claude-code-9dc16082](../actor/claude-code-9dc16082.md)

## result

PASS

## statement

QA at b9129e145 with a scratch build in a scratch root and project: rhiza daemon settings lists runawayProcs 400, runawayMemory 8G, runawayGrace 5m (default) for machine and project; daemon set machine runawayMemory 12G/512M and runawayProcs 0 land in rhiza.toml, project runawayProcs 50/runawayMemory 2G/runawayGrace 90s land in git config and show (git config, project); 'runawayMemory lots' refused naming the key (exit 2), '-- runawayProcs -1' and '1.5' refused naming the key, runawayGrace 0 refused. go test ./internal/api -run TestWatchdogPolicy (runaway_keys, project_over_machine) and TestScripts/watchdog-policy exit 0; rhiza docs --out scratch gives config.md byte-identical to docs/reference/config.md.

## Links

### Out

- verifies [runaway-threshold-keys](../capability/rhiza/orchestration/daemon/daemon-watchdog/runaway-threshold-keys.md)

## History

- `01M7NVMDMBVBKRW5B8SB0KMW6T` The daemon's watchdog treats a runaway agent session of any workflow, one whose process tree's process count or resident memory grows past a host threshold read portably (ps, on Linux and macOS), like a stall: it nudges the session, then fails the step with the reading as its reason
