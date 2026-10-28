# daemon-run-tabs-verified

evidence · `01M628XSKNAS3HT1558KFTHR3D` · actor [claude-code-39aa7c66](../actor/claude-code-39aa7c66.md)

## result

PASS

## statement

At 5a5d694, go test ./cmd/rhiza -run 'TestScripts/daemon-picks$' -testwork passes; its fake Rex logged, for rhiza daemon --once, one new tab 'rex window new 01M5PWGQ8 -c <repo> --json -- <rhiza> run start 01M5PWGQ8GB4618G4RH2TV0ZAH' (the Built change set's ID, no --autonomous); the script asserts one such line and zero --autonomous before the flag. rhiza daemon --once --autonomous logged the same command ending in --autonomous, asserted by the script's grep (exactly one line). Diff: DaemonRequest gains Autonomous, threaded through daemonPoll1 to daemonStart, which appends --autonomous only when set; describe-golden updated. make fmt-check, vet and test pass (self-check's only failures are this verifier session's own actor render). The script pins the ID only via the logged line, not its own assertion (regex [0-9A-Z]*).

## Links

### Out

- verifies [daemon-run-tabs](../capability/rhiza/orchestration/daemon/daemon-run-tabs.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
