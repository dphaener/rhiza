# ci-unstable-repeated-flake-verified

evidence · `01M78X24PJWRJJPJJAFTJWNN78` · actor [claude-code-57f08c13](../actor/claude-code-57f08c13.md)

## result

PASS

## statement

Verifier at 3cb4a4ff5: go test -trimpath -count=1 ./internal/api -run 'TestCIStability|TestShipLocalGate' -v passes (TestCIStability, TestShipLocalGate, TestShipLocalGateOncePerHead). The cases cover: one TestX flake gives stable with reason 'TestX was seen once and ignored'; the same head twice gives stable; TestX and TestY once each give stable; TestX on 2 heads gives unstable with reason 'TestX flaked 2 times on 2 heads in the last 20 CI runs (stops applying after 16 more completed CI runs)', with N=2 for window 5, and stable after exactly N more clean runs but not N-1; infra-only gives stable with reason 'ignored infra failures: ...'; red main gives unstable with no flake and is cleared by a ship check main green record. Read the diff: the N formula window-(len-idx)+1 matches when the older of the two latest flaky heads leaves the window. ship_gate.go appends any stable verdict reason to 'CI is stable (...)', and the infra-only gate test checks that path, so a single-flake skip carries 'seen once and ignored'. ciVerdict takes its window from rhiza.ship.ciWindow through ciStabilityVerdict. gofmt and go vet on internal/api are clean.

## Links

### Out

- verifies [ci-unstable-repeated-flake](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-unstable-repeated-flake.md)

## History

- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
