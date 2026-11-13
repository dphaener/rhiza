# rhiza-scoped-gate-verified-2

evidence · `01M7936C4D2PX5K2ZH486XW5RG` · actor [claude-code-ffbb4bf3](../actor/claude-code-ffbb4bf3.md)

## result

PASS

## statement

Verified at 53b6f13a9 (fix for quality finding f1). go test -trimpath -count=1 -v ./tools/gate/... passes: TestPlanOnePackageAndDependents (internal/a change plans go test ./internal/a ./internal/b plus cmd/rhiza -skip, no internal/c, no make test, then make self-check), TestPlanRhizaOnly (.rhiza/-only change plans make self-check only), TestPlanFallbacks (empty RHIZA_GATE_BASE, go.mod, and now docs/quickstart.md, AGENTS.md, .github/workflows/ci.yml each plan make test plus make self-check), TestPlanScripts (changed foo.txtar with acceptance naming TestScripts/bar plans -run 'TestScripts/(bar|foo)$', no internal/), TestScriptNames. grep -n "rhiza.ship.localGate 'make gate'" AGENTS.md finds line 96. Re-ran the f1 repro RHIZA_GATE_BASE=3980a36aa~1 RHIZA_GATE_HEAD=3980a36aa go run ./tools/gate -dry-run: now plans make test then make self-check (README.md outside every package) instead of self-check only. Read the diff of plan.go (non-.rhiza/ files outside packages return fullSteps) and AGENTS.md wording. gofmt -l and go vet ./tools/gate clean.

## Links

### Out

- verifies [rhiza-scoped-gate](../capability/rhiza/integrity/self-check/rhiza-scoped-gate.md)

## History

- `01M78T0G869AE4NQ9YGNB2DDGE` A passing local gate is reused while the code it checked hasn't changed, and the gate checks only what the change touches. Owner, 2026-11-12, agreed to file after the performance report. Today run 01M787SV ran the full make check 4 times (local-gate, gate-rebase-merge, after conflict-fix, after a rendered-only merge), ~20 min each serialized on the test lock; most reruns followed changes that touched no Go code (rendered model files, .rhiza/log entries). Wanted: (1) record the tree hash (or the hash of the non-.rhiza/ tree) a gate passed on and skip the gate when the next push has the same code tree; (2) the gate runs the change's scoped checks (go test for the packages the branch changes vs origin/main, their reverse dependencies, the affected testscripts, make self-check) instead of the whole suite, leaving the full suite to CI; (3) the project setting stays language-neutral (the gate command and its scoping come from project/workflow config, per the cleanup milestone). Estimated saving ~20 min per run on days the gate is on. Evidence: ~/orchestrator/perf-2026-11-12/report-prospective.md.
