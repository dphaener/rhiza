# modellint-unapplied-decisions

capability · `01M7AZSBANC9KR2Y0J1QZZTWZ9` · parent [built-in-review](../built-in-review.md)

## acceptance

go test -trimpath ./internal/modellint/... passes with table tests in internal/modellint/decisions_test.go: a Decision no node links decided-by to is reported; an answers-only Decision is not; one linked from a landed capability is not; one linked only from leaves of the change set none of which has a Leaf: commit is reported, and with one such commit given, by handle or by ID, it is not; failure case: an unapplied Decision another change set wrote is not reported.

## statement

internal/modellint reports each Decision the change set adds or changes that is never applied: no live node links to it with decided-by (a Decision whose only link is answers on a question is exempt from this clause), or every node that does is a leaf of this change set with no branch commit carrying its Leaf: trailer, the leaves with such commits given to it as input by handle or ID. It judges nothing of the code's meaning and changes no file outside internal/modellint.

## Links

### Out

- depends-on [modellint-stale-refs](modellint-stale-refs.md)
- decided-by [lint-hits-feed-model-check](../../../../../decision/lint-hits-feed-model-check.md)

### In

- depends-on [code-check-runs-lints](code-check-runs-lints.md)
- verifies [modellint-unapplied-decisions-verified](../../../../../evidence/modellint-unapplied-decisions-verified.md)
- verifies [modellint-unapplied-decisions-qa-1](../../../../../evidence/modellint-unapplied-decisions-qa-1.md)
- verifies [pr-569-model-lints-in-review](../../../../../evidence/pr-569-model-lints-in-review.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
