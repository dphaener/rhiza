# rhiza-scoped-gate-qa-1

evidence · `01M792XN2Y1N4VB07K0AQHPG85` · actor [claude-code-7ac56481](../actor/claude-code-7ac56481.md)

## result

PASS

## statement

QA in a detached worktree of the branch head: make gate ARGS=-dry-run printed make test + make self-check for an empty RHIZA_GATE_BASE ('no base to scope by') and for the merge-base with origin/main ('Makefile changed'); make self-check only for base 8e79a27 (6 .rhiza/ files changed); after a probe commit to internal/textref, go test -trimpath for textref and its 13 transitive reverse dependencies, cmd/rhiza non-script tests and TestScripts limited to the leaf-named scripts under the test lock, then make self-check; after a probe change to ship-local-gate.txtar with RHIZA_CHANGE set, the -run pattern included ship-local-gate and describe-golden. A real make gate with base 8e79a27 ran make self-check only and exited 0 (with the run's RHIZA_ROOT unset; under it, self-check fails on the run's own unrendered actors, which is the run environment, not the change). go test -trimpath ./tools/gate passes; AGENTS.md:95 has the localGate 'make gate' line.

## Links

### Out

- verifies [rhiza-scoped-gate](../capability/rhiza/integrity/self-check/rhiza-scoped-gate.md)

## History

- `01M78T0G869AE4NQ9YGNB2DDGE` A passing local gate is reused while the code it checked hasn't changed, and the gate checks only what the change touches. Owner, 2026-11-12, agreed to file after the performance report. Today run 01M787SV ran the full make check 4 times (local-gate, gate-rebase-merge, after conflict-fix, after a rendered-only merge), ~20 min each serialized on the test lock; most reruns followed changes that touched no Go code (rendered model files, .rhiza/log entries). Wanted: (1) record the tree hash (or the hash of the non-.rhiza/ tree) a gate passed on and skip the gate when the next push has the same code tree; (2) the gate runs the change's scoped checks (go test for the packages the branch changes vs origin/main, their reverse dependencies, the affected testscripts, make self-check) instead of the whole suite, leaving the full suite to CI; (3) the project setting stays language-neutral (the gate command and its scoping come from project/workflow config, per the cleanup milestone). Estimated saving ~20 min per run on days the gate is on. Evidence: ~/orchestrator/perf-2026-11-12/report-prospective.md.
