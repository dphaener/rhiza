# ship-gate-full-log-qa-1

evidence · `01M78NBNR81BYAVENYRBDRY6HW` · actor [claude-code-6bfb93b7](../actor/claude-code-6bfb93b7.md)

## result

PASS

## statement

QA by hand with the built binary (go build ./cmd/rhiza at 952263bc7) in a fresh scratch repo (git init, rhiza init, change new, rhiza daemon set rhiza.ship.localGate to a script, a ci-stability.json with an identified flake so CI is unstable). Failing gate printing '--- FAIL: TestEarly' and 'FAIL example.com/pkg/early' first, 30000 filler lines (619039 bytes), then make's 'Error 1', exit 2: rhiza ship gate exited 1; its last 20 human lines (a step's kept tail) held '--- FAIL: TestEarly', the FAIL line, 'Full log: <common>/rhiza/gate/cs/qa-gate-log.log' and the failure line; --json log names that file, output kept at the 16 KB tail; the log file is byte-identical to the gate's full output (cmp). A passing gate (echo all fine) exited 0 and replaced the log with 'all fine'; rerun on the same gated head passes from the cache without rerunning or naming a log. grep -c 'run `rhiza ship gate` to see' prints 0 for ship-change.toml and docs.toml; both gate-fix prompts in each name `rhiza run status --json` and the full gate log file. go test -trimpath ./internal/api -run 'TestShipGateFullLog|TestShipRebase|TestShipLocalGate' and ./cmd/rhiza TestScripts/(describe-golden|ship-local-gate) pass. Rebase-path log write checked through its unit test only.

## Links

### Out

- verifies ship-gate-full-log (retired)

## History

- `01M78FD039D0F1QVN2BWEAEYCC` Fix The local ship gate (rhiza.ship.localGate = make check) fails in every run under the VM root: internal/api passes all tests but its leak guard (internal/rootguard) ends 'FAIL: the test run wrote into the state folder of the root enclosing it: .rhiza-root/usage.live.json: file changed', because every live Claude session's status line rewrites ~/code/.rhiza-root/usage.live.json during the ~10-20 min suite. Evidence 2026-11-12 16:30-17:10: runs 01M787SV, 01M788HB and 01M7859E all failed local-gate/gate-rebase-merge at 'make: *** [Makefile:64: test] Error 1' and went to gate-fix; the gate was on only because one 'scripts (4/4) / UNKNOWN STEP' flake on #498 marked CI unstable (29 of the last 30 CI runs green). The gate-fix for 01M7859E fixed the guard inside an unrelated change set (commit e075aefa6 'Judge the live usage file by captured_at in the root guard', PR #498 rhiza release check); the other two fixers were working the same failure in parallel, so expect duplicate fixes or conflicts. Also: the gate output stored in the journal is only the last ~2 KB, which cut off the failing test's name, so every fixer had to rerun the full suite (~20 min each, serialized on the test lock). Fix ideas: land the guard fix on its own; store or point to the full gate log; consider whether one unidentified flake should switch on the full local gate for every run.
