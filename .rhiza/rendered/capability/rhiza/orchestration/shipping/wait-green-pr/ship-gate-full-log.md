# ship-gate-full-log

capability · `01M78FP45EF4RB94DP1JZ5W6T2` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/api -run 'TestShipGateFullLog|TestShipRebase' shows a gate whose output is over 64 KB with a '--- FAIL: TestX' line near its start and make's error lines at its end: the log file holds every byte, ShipGateResponse's log field names it, and the last 20 lines of the human output (what a failed command step keeps) contain '--- FAIL: TestX' and the log path; a passing gate also writes the log, and a failing gate run by ship rebase writes the same log; describe-golden is regenerated for the new field; grep -c 'run `rhiza ship gate` to see' .rhiza/workflows/ship-change.toml .rhiza/workflows/docs.toml prints 0 for both, each gate-fix prompt names `rhiza run status --json` and the full log, and make self-check passes

## statement

When the local gate runs (rhiza ship gate, and the gate ship rebase runs on a merged head, both through shipLocalGate), its whole combined output is written to a log file under the git common dir (rhiza/gate/<branch>.log, replaced on each run), the response carries the file's path, and a failing gate's human output ends with the output's failing test lines (go test '--- FAIL:', 'FAIL', 'panic:' lines, at most 10) and 'Full log: <path>', so the tail a run step keeps names the failing test and where to read the rest. The gate-fix prompts (.rhiza/workflows/ship-change.toml and docs.toml, two each) stop telling the fixer to run `rhiza ship gate` to see the gate's output, since a failed head is never cached and that reruns the whole gate; they tell it to read the failed gate step's kept output in `rhiza run status --json` and the full log file it names

## Links

### Out

- depends-on [ship-local-gate](ship-local-gate.md)
- decided-by [gate-log-kept-whole](../../../../../decision/gate-log-kept-whole.md)

### In

- verifies [ship-gate-full-log-verified](../../../../../evidence/ship-gate-full-log-verified.md)
- verifies [ship-gate-full-log-qa-1](../../../../../evidence/ship-gate-full-log-qa-1.md)
- verifies [pr-504-fix-the-local-ship-gate](../../../../../evidence/pr-504-fix-the-local-ship-gate.md)

## History

- `01M78FD039D0F1QVN2BWEAEYCC` Fix The local ship gate (rhiza.ship.localGate = make check) fails in every run under the VM root: internal/api passes all tests but its leak guard (internal/rootguard) ends 'FAIL: the test run wrote into the state folder of the root enclosing it: .rhiza-root/usage.live.json: file changed', because every live Claude session's status line rewrites ~/code/.rhiza-root/usage.live.json during the ~10-20 min suite. Evidence 2026-11-12 16:30-17:10: runs 01M787SV, 01M788HB and 01M7859E all failed local-gate/gate-rebase-merge at 'make: *** [Makefile:64: test] Error 1' and went to gate-fix; the gate was on only because one 'scripts (4/4) / UNKNOWN STEP' flake on #498 marked CI unstable (29 of the last 30 CI runs green). The gate-fix for 01M7859E fixed the guard inside an unrelated change set (commit e075aefa6 'Judge the live usage file by captured_at in the root guard', PR #498 rhiza release check); the other two fixers were working the same failure in parallel, so expect duplicate fixes or conflicts. Also: the gate output stored in the journal is only the last ~2 KB, which cut off the failing test's name, so every fixer had to rerun the full suite (~20 min each, serialized on the test lock). Fix ideas: land the guard fix on its own; store or point to the full gate log; consider whether one unidentified flake should switch on the full local gate for every run.
