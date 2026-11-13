# rhiza-scoped-gate

capability · `01M78TAJENDS93MS29H076VKXT` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath ./tools/gate shows the plan the tool prints with -dry-run in a temporary git repo: a change to one internal package lists go test for that package and its reverse dependencies only, then make self-check; a change only under .rhiza/ lists make self-check only; a change to go.mod, or an empty RHIZA_GATE_BASE, lists make test and make self-check; a changed cmd/rhiza/testdata/script/foo.txtar with a leaf acceptance naming TestScripts/bar lists cmd/rhiza with -run 'TestScripts/(bar|foo)$'; grep -n "rhiza.ship.localGate 'make gate'" AGENTS.md finds the line

## statement

Rhiza's own repo has a scoped gate, make gate running a Go tool under tools/gate, that from RHIZA_GATE_BASE runs go test -trimpath on the Go packages the branch changes and their reverse dependencies (cmd/rhiza's TestScripts limited to the changed txtar files and those named TestScripts/<name> in the leaf acceptance checks of RHIZA_CHANGE), then make self-check, with cmd/rhiza under the machine-wide test lock, and falls back to make test plus make self-check when the base is empty or go.mod, go.sum, the Makefile, tools/gate or cmd/rhiza/main_test.go changed; AGENTS.md says to set rhiza.ship.localGate to 'make gate'

## Links

### Out

- depends-on [ship-gate-base-env](../../orchestration/shipping/wait-green-pr/ship-gate-base-env.md)
- decided-by [local-gate-reuse-and-scope](../../../../decision/local-gate-reuse-and-scope.md)

### In

- verifies [rhiza-scoped-gate-verified](../../../../evidence/rhiza-scoped-gate-verified.md)
- verifies [rhiza-scoped-gate-qa-1](../../../../evidence/rhiza-scoped-gate-qa-1.md)
- verifies [rhiza-scoped-gate-verified-2](../../../../evidence/rhiza-scoped-gate-verified-2.md)
- verifies [pr-513-a-passing-local-gate-is-reused](../../../../evidence/pr-513-a-passing-local-gate-is-reused.md)

## History

- `01M78T0G869AE4NQ9YGNB2DDGE` A passing local gate is reused while the code it checked hasn't changed, and the gate checks only what the change touches. Owner, 2026-11-12, agreed to file after the performance report. Today run 01M787SV ran the full make check 4 times (local-gate, gate-rebase-merge, after conflict-fix, after a rendered-only merge), ~20 min each serialized on the test lock; most reruns followed changes that touched no Go code (rendered model files, .rhiza/log entries). Wanted: (1) record the tree hash (or the hash of the non-.rhiza/ tree) a gate passed on and skip the gate when the next push has the same code tree; (2) the gate runs the change's scoped checks (go test for the packages the branch changes vs origin/main, their reverse dependencies, the affected testscripts, make self-check) instead of the whole suite, leaving the full suite to CI; (3) the project setting stays language-neutral (the gate command and its scoping come from project/workflow config, per the cleanup milestone). Estimated saving ~20 min per run on days the gate is on. Evidence: ~/orchestrator/perf-2026-11-12/report-prospective.md.
