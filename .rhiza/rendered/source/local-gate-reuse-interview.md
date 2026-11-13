# local-gate-reuse-interview

source · `01M78TA0E1BBZE2E039NJ4Z97N`

## quotes

```json
[]
```

## raw_sha256

3e67215be090a992503541150e7399a11a6bc56a69eb98df410ee20e93b61884

## session

aea6d8fc-7b79-4c73-8624-01a1d1a09694

## summary

Autonomous plan run (no owner) for change set 01M78T0, filed after the 2026-11-12 performance report (~/orchestrator/perf-2026-11-12/report-prospective.md, recommendation 1: run 01M787SV ran the full make check 4 times, ~20-25 min each on the test lock, mostly after changes that touched only .rhiza/). The planner read internal/api/ship_gate.go (the gate already passes at once on a head recorded as Gated in the ship state, keyed on the commit, so any render or log commit reruns it), the ship-change workflows and the sibling change set 01M78RM (gate switches on only for repeated named flakes). Three questions went to the arbiter in one call: which intent to serve, what the reuse key is, and how scoping stays language-neutral. The arbiter's Decision local-gate-reuse-and-scope: serve m35-ship-40m; record the hash of the passing head's tree with top-level .rhiza/ left out and pass at once when a clean HEAD has the same non-.rhiza/ tree (.rhiza/ is checked on every push by the pre-push check); the gate command stays the project setting rhiza.ship.localGate, run with RHIZA_GATE_BASE (merge-base with origin/main) and RHIZA_GATE_HEAD in its environment, and Rhiza's own repo supplies a script that scopes go test to changed packages, their reverse dependencies, the affected testscripts and make self-check, falling back to the full suite; no Go-specific scoping in the shipped product. The planner chose the affected testscripts as the changed txtar files plus those named TestScripts/<name> in the change set's leaf acceptance checks, and a reuse that also requires the same gate command.

## title

Plan of 01M78T0: a passing local gate is reused

## Links

### In

- derived-from [local-gate-reuse-and-scope](../decision/local-gate-reuse-and-scope.md)

## History

- `01M78T0G869AE4NQ9YGNB2DDGE` A passing local gate is reused while the code it checked hasn't changed, and the gate checks only what the change touches. Owner, 2026-11-12, agreed to file after the performance report. Today run 01M787SV ran the full make check 4 times (local-gate, gate-rebase-merge, after conflict-fix, after a rendered-only merge), ~20 min each serialized on the test lock; most reruns followed changes that touched no Go code (rendered model files, .rhiza/log entries). Wanted: (1) record the tree hash (or the hash of the non-.rhiza/ tree) a gate passed on and skip the gate when the next push has the same code tree; (2) the gate runs the change's scoped checks (go test for the packages the branch changes vs origin/main, their reverse dependencies, the affected testscripts, make self-check) instead of the whole suite, leaving the full suite to CI; (3) the project setting stays language-neutral (the gate command and its scoping come from project/workflow config, per the cleanup milestone). Estimated saving ~20 min per run on days the gate is on. Evidence: ~/orchestrator/perf-2026-11-12/report-prospective.md.
