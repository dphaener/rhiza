# run-start-new-intent-verified

evidence · `01M6A6K92CVJNDSS8NZFBJ5D26` · actor [claude-code-007d2717](../actor/claude-code-007d2717.md)

## result

PASS

## statement

At 0bb2f0a6 I ran go test ./cmd/rhiza -run 'TestScripts/(run-start-new|describe-golden)$' -count=1 (ok) and make check with the run's RHIZA_ROOT/CHANGE/RUN unset (fmt, vet, all tests and self-check pass). run-start-new.txtar uses fake rex and claude binaries and shows: a change set argument exits 2 with 'takes no change set'; --workflow other exits 2 with 'not plan'; neither opens a change set or a run; --dry-run prints a change ID and writes nothing; the real run succeeds, and change ls gains a draft intention whose id is the predicted one, or a draft record with --milestone; fold-ok passes. I also checked by hand with a scratch build, fake rex/claude and a fixed clock and seed in a repo that already has a change set: the dry run predicted 01M5PWGJC82P7YGZTWYFFZ0VJG, the real run opened that same ID as a draft intention, and the run journal's change field named it. I read the diff in internal/api/run.go (newPlanChange, the --milestone without --new usage error) and the describe golden update. Side note: in a fresh repo with no recorded actor, the dry-run prediction misses because the real run first records the actor. Plain change new --dry-run does the same, so this leaf did not cause it.

## Links

### Out

- verifies [run-start-new-intent](../capability/rhiza/roadmap-app-start-runs/run-start-new-intent.md)

## History

- `01M686RVAD1EHGERBBDJN3S45G` The roadmap app starts a plan run for a new milestone or a new change set from a typed one-line intent
