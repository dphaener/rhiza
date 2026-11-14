# code-check-runs-lints

capability · `01M7AZSBEY65G73S64CERYNTD0` · parent [built-in-review](../built-in-review.md)

## acceptance

go test -trimpath ./cmd/rhiza -run 'TestScripts/(review-code-check|describe-golden)$' passes: in review-code-check.txtar a change set whose leaf acceptance names a missing file and which adds a Decision nothing links decided-by to, on a branch that deletes the only file a landed capability names, gets code-check exit 0 with all three hits under lints in --json and under code_checks in rhiza run status --json; failure case: the same change set with a leaf lacking verifier Evidence exits 1 and its causes name only the Evidence, not the lints. go test -trimpath ./internal/run/... ./internal/api/... passes with the code-check entry's lints round-tripping through the journal; describe-golden.txtar is regenerated with -update.

## statement

`rhiza run code-check` runs the internal/modellint lints (stale references, unapplied Decisions and capabilities whose code the branch removed) on the run's change set in the run's worktree (the head tree, the Leaf: trailers of the commits from the merge-base to HEAD and the paths deleted in merge-base..HEAD) and lists their hits as `lints` in its --json response, in its human output and on the journalled code-check entry, which `rhiza run status --json` shows under code_checks; a lint hit is never a cause and never changes `passed` or the exit code. Files: internal/api/run_codecheck.go, internal/run/journal.go, internal/api/run.go, cmd/rhiza/testdata/script/review-code-check.txtar, cmd/rhiza/testdata/script/describe-golden.txtar.

## Links

### Out

- depends-on [modellint-unapplied-decisions](modellint-unapplied-decisions.md)
- depends-on [modellint-branch-removed-code](../../../../modellint-branch-removed-code.md)
- decided-by [lint-hits-feed-model-check](../../../../../decision/lint-hits-feed-model-check.md)

### In

- depends-on [review-lint-fix-route](../../../../review-lint-fix-route.md)
- verifies [code-check-runs-lints-verified](../../../../../evidence/code-check-runs-lints-verified.md)
- verifies [code-check-runs-lints-qa-1](../../../../../evidence/code-check-runs-lints-qa-1.md)
- verifies [pr-569-model-lints-in-review](../../../../../evidence/pr-569-model-lints-in-review.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
