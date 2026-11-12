# health-ci-infra-failures-verified

evidence · `01M78X9EW1Z40YB2MRENXAFDSN` · actor [claude-code-2a609595](../actor/claude-code-2a609595.md)

## result

PASS

## statement

At fa3d982ee: go test -trimpath ./internal/api -run TestHealthCI -v passes. In TestHealthCI, a record with one 'scripts (4/4) / UNKNOWN STEP' infra failure and one TestX flake gives ok/stable. Its detail lists 'flaky TestX: 1 flake' and 'infra scripts (4/4) / UNKNOWN STEP: 1 failure ... infra failures do not make CI unstable'. The stable summary comes from ciVerdict, which reads 'TestX was seen once and ignored; ignored infra failures: ...'. A record where TestX flaked on 2 heads gives warn, naming 'TestX flaked 2 times on 2 heads' and 'stops applying after'. TestHealthCIStability also passes. go test -trimpath ./cmd/rhiza -run TestScripts/health-ci passes, and its script checks the infra line in both rhiza health and rhiza health --json. I read the diff: ciFailureCounts counts infra per job / step in the window, and healthCI appends one sorted infra detail line for each.

## Links

### Out

- verifies [health-ci-infra-failures](../capability/rhiza/orchestration/shipping/wait-green-pr/health-ci-infra-failures.md)

## History

- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
