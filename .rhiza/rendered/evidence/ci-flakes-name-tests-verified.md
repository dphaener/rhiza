# ci-flakes-name-tests-verified

evidence · `01M78WGPY9C2PHMNW5HJ9KKB7G` · actor [claude-code-ff27ba14](../actor/claude-code-ff27ba14.md)

## result

PASS

## statement

Verifier at ddcc0cf53, detached checkout. Ran go test -trimpath ./internal/api -run 'TestShipCheckStabilityRecord|TestFlakyTestBugs|TestTrimFailureLog' -v: all PASS. Subtests cover each acceptance point: 'subtests by full name' (a log with '--- FAIL: TestFoo' and an indented '--- FAIL: TestScripts/bar', rerun green) records Flakes [TestFoo TestScripts/bar] and no Infra; 'unknown step' records Infra ['scripts (4/4) / UNKNOWN STEP'] and no flake; 'no job or step' records Infra ['CI run 42']; 'red after rerun' records red with no flakes (the code sets Flakes/Infra only on the green branch of recordSettledCI, so Infra is empty too); 'old record reads job / step flakes as infra' splits ['test / Run','CI run 9','TestA'] into Flakes [TestA] and Infra on read; TestFlakyTestBugs queues bugs for named test flakes and none for an infra-only run (Infra entry or a legacy 'CI run 5' flake). Also ran TestCIStability and TestHealthCIStability (PASS), TestScripts/describe-golden (ok), go vet and gofmt -l on internal/api (clean). Read the diff: failedTestNames in ship.go matches '--- FAIL: <name>' after logText strips the job/step prefix, timestamp and indentation, and drops a parent whose subtest is named; shipFailedTests sends a log that names no test to infra; CIRun.normalized moves 'job / step' and 'CI run' flakes into Infra on read, and ciVerdict, ciFlakeCounts and watchdogFlakyTests all use it.

## Links

### Out

- verifies [ci-flakes-name-tests](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-flakes-name-tests.md)

## History

- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
