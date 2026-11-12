# local-gate-switches-plan

source · `01M78S68PFA2DSYGZSEPYHAWBD`

## quotes

```json
[]
```

## raw_sha256

025b3a3d920301892b08ca2cb62d10d6b60fb27de2875004d850b82d9cee39b1

## session

250a1536-7da7-4f36-b513-01a1d18ad9d4

## summary

Autonomous plan run (no owner) for change set 01M78RM, filed after the 2026-11-12 performance report: one 'scripts (4/4) / UNKNOWN STEP' flake marked CI unstable and switched the ~20 min serialized local gate on for every run (~525 run-minutes lost), and the journal kept only the gate's last ~2 KB, cutting off the failing test. Read ci_stability.go (any flake in the window makes CI unstable), ship.go (flakes are recorded as 'job / step' or 'CI run <id>', never a test name), ship_gate.go (output kept in memory, 16 KB tail), run.go (a failed step keeps 2 KB / 20 lines), watchdog_flaky.go and the gate-fix prompts (they tell the fixer to rerun rhiza ship gate). Related accepted bug 01M78FD0 already plans a full gate log (ship-gate-full-log) and ignoring unidentified flakes. One question went to the arbiter, who ruled (local-gate-repeated-named-flakes): serve m4-runs-hold-up; a failure naming no test is infra, recorded apart, shown in health, never unstable, no bug filed; a named test is unstable only after flaking on 2 distinct heads in the ciWindow, with a reason naming it, its count and when it stops applying; each gate run writes its whole output to a log file in the run dir, its path in the response and the journal, the failure output ending with the failing test lines and the path; gate-fix prompts read the log. Planned four leaves: named flakes vs infra in the record, the repeated-flake verdict, infra in health, and the per-run gate log. Because 01M78FD0's nodes are not landed, the leaves are new capabilities and the change set depends on 01M78FD0.

## title

Local ship gate switches, autonomous plan

## Links

### In

- derived-from [local-gate-repeated-named-flakes](../decision/local-gate-repeated-named-flakes.md)

## History

- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
