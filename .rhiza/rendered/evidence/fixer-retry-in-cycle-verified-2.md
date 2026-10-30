# fixer-retry-in-cycle-verified-2

evidence · `01M658NNXWQJ0AECY3GZ09K1QA` · actor [claude-code-c97be403](../actor/claude-code-c97be403.md)

## result

PASS

## statement

At a7ffc29 (attempt 2, a model-only fix that narrows decision fixer-retries-within-cycle to the change set's fixer, with the leaf's acceptance unchanged), I ran go test ./internal/run -run TestFixChangeSet -count=1 -v. All of them passed, including TestFixChangeSetRetriedInCycle: the change-set fixer fails attempt 1 and passes attempt 2, the run ends review:pass with 2 cycles (no new cycle), the fixer specs are cycle 2 attempts 1 and 2, and the journal's 'review-fix-judged fixer 2 pass' comes before 'check-session alpha check 2'. TestFixChangeSetActions 'its failure is retried in the cycle' also passed. I read changeSetFix in internal/run/checks.go: a judged-failed fixer with Attempt <= MaxFixCycles and no leaf failure is started again with attempt+1 in the same cycle. go test ./... -count=1 passed. go vet passed and gofmt -l was clean. make self-check failed only on stale .rhiza/rendered files for actors added by the orchestrator's check-write commits 8c751c9 and cc58773, not by this leaf; the run's render step covers those. One uncached go test ./... run hit the known TestGeneratedModel timing flake under load, and it passed on the -count=1 run.

## Links

### Out

- verifies [fixer-retry-in-cycle](../capability/fixer-retry-in-cycle.md)

## History

- `01M61GW7TJ6BH6YX1T7SQDWK1J` A failed fixer does not use up a review fix cycle: the fixer is retried within its cycle, and the cycle's checks re-run only once its fixers pass.
