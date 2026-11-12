# ship-gate-full-log-verified

evidence · `01M78N3RXWFN2GNV6HX2V6181K` · actor [claude-code-4658e5c7](../actor/claude-code-4658e5c7.md)

## result

PASS

## statement

Verified at eee7eb49d (attempt 2). go test -trimpath ./internal/api -run 'TestShipGateFullLog|TestShipRebase' -count=1: all PASS. TestShipGateFullLog runs a gate printing '--- FAIL: TestX' first, 20000 filler lines (>64 KB) and make's 'Error 1' last: the log file holds every byte, ShipGateResponse.Log names it, and the last 20 human lines contain '--- FAIL: TestX' and the log path; its passing subtest shows a passing gate writes the log too. Attempt 1's finding is fixed: TestShipRebaseRunsLocalGate's failing case now asserts r.Gate.Log is <common>/rhiza/gate/cs/ship-it.log, the file holds the whole output, and the human output has 'Full log: <path>' (shared path through shipLocalGate, ship_rebase.go). Diff read: shipWriteGateLog writes atomically under the common dir, shipGateFailLines keeps at most 10 --- FAIL/FAIL/panic: lines printed before 'Full log:'. describe-golden has the log field and TestScripts/describe-golden passes. grep -c 'run `rhiza ship gate` to see' prints 0 for ship-change.toml and docs.toml; both gate-fix prompts in each name `rhiza run status --json` and the full gate log. make self-check passes (exit 0) from a clean detached worktree of eee7eb49d with the run's RHIZA_* env unset (with it set, it fails only on the verifier sessions' own unrendered actor files).

## Links

### Out

- verifies ship-gate-full-log (retired)

## History

- `01M78FD039D0F1QVN2BWEAEYCC` Fix The local ship gate (rhiza.ship.localGate = make check) fails in every run under the VM root: internal/api passes all tests but its leak guard (internal/rootguard) ends 'FAIL: the test run wrote into the state folder of the root enclosing it: .rhiza-root/usage.live.json: file changed', because every live Claude session's status line rewrites ~/code/.rhiza-root/usage.live.json during the ~10-20 min suite. Evidence 2026-11-12 16:30-17:10: runs 01M787SV, 01M788HB and 01M7859E all failed local-gate/gate-rebase-merge at 'make: *** [Makefile:64: test] Error 1' and went to gate-fix; the gate was on only because one 'scripts (4/4) / UNKNOWN STEP' flake on #498 marked CI unstable (29 of the last 30 CI runs green). The gate-fix for 01M7859E fixed the guard inside an unrelated change set (commit e075aefa6 'Judge the live usage file by captured_at in the root guard', PR #498 rhiza release check); the other two fixers were working the same failure in parallel, so expect duplicate fixes or conflicts. Also: the gate output stored in the journal is only the last ~2 KB, which cut off the failing test's name, so every fixer had to rerun the full suite (~20 min each, serialized on the test lock). Fix ideas: land the guard fix on its own; store or point to the full gate log; consider whether one unidentified flake should switch on the full local gate for every run.
