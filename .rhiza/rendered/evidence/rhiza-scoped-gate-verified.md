# rhiza-scoped-gate-verified

evidence · `01M792GEPG2SJZAEGTRJD5BTG4` · actor [claude-code-dcfb13e5](../actor/claude-code-dcfb13e5.md)

## result

PASS

## statement

Verifier at 8e79a2706: go test -trimpath -count=1 -v ./tools/gate passed all five tests. TestPlanOnePackageAndDependents runs -dry-run in a temp git repo and shows that a change to internal/a lists go test ./internal/a ./internal/b (its reverse dependency) plus cmd/rhiza, not internal/c or make test, then make self-check. TestPlanRhizaOnly shows a .rhiza/-only change lists make self-check alone. TestPlanFallbacks shows that an empty RHIZA_GATE_BASE and a go.mod change each list make test and make self-check. TestPlanScripts shows a changed foo.txtar plus an acceptance naming TestScripts/bar lists -run 'TestScripts/(bar|foo)$'. grep -n "rhiza.ship.localGate 'make gate'" AGENTS.md matched lines 95-96. I also dry-ran the tool on this repo. Base HEAD~1 shows a Makefile change, so it falls back to make test plus make self-check. HEAD~3..HEAD~2 (internal/api plus describe-golden.txtar) lists internal/api with its reverse dependencies, then cmd/rhiza under the test lock with -skip '^TestScripts$', then TestScripts/(bar|describe-golden|foo|ship-local-gate)$, then make self-check. Minor observation: bar and foo come from this leaf's own example acceptance text and match no scripts, which does no harm. A .rhiza-only commit lists make self-check only. The diff matches the leaf: make gate runs go run ./tools/gate, tools/gate changes trigger the fallback, and cmd/rhiza steps run under the flock test lock. go vet and gofmt are clean.

## Links

### Out

- verifies [rhiza-scoped-gate](../capability/rhiza/integrity/self-check/rhiza-scoped-gate.md)

## History

- `01M78T0G869AE4NQ9YGNB2DDGE` A passing local gate is reused while the code it checked hasn't changed, and the gate checks only what the change touches. Owner, 2026-11-12, agreed to file after the performance report. Today run 01M787SV ran the full make check 4 times (local-gate, gate-rebase-merge, after conflict-fix, after a rendered-only merge), ~20 min each serialized on the test lock; most reruns followed changes that touched no Go code (rendered model files, .rhiza/log entries). Wanted: (1) record the tree hash (or the hash of the non-.rhiza/ tree) a gate passed on and skip the gate when the next push has the same code tree; (2) the gate runs the change's scoped checks (go test for the packages the branch changes vs origin/main, their reverse dependencies, the affected testscripts, make self-check) instead of the whole suite, leaving the full suite to CI; (3) the project setting stays language-neutral (the gate command and its scoping come from project/workflow config, per the cleanup milestone). Estimated saving ~20 min per run on days the gate is on. Evidence: ~/orchestrator/perf-2026-11-12/report-prospective.md.
