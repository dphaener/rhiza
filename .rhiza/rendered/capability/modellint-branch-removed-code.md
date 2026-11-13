# modellint-branch-removed-code

capability · `01M7B0BDZ5TW0VHH9KH1XMR3E1`

## acceptance

go test -trimpath ./internal/modellint/... passes with table tests in internal/modellint/removed_test.go: a landed capability whose only path the deleted list holds and the head tree lacks is reported naming the capability and the path; failure cases: one with a deleted path and another path still present is not reported, one whose path is missing but not in the deleted list (removed before the branch) is not reported, a retired capability is not reported, and a capability with no path references is not reported.

## statement

internal/modellint reports as unretired each live capability outside the change set's slice (any landed or in-flight capability the change set does not write) that references at least one repo path, whose every path reference is absent from the head tree, and at least one of whose path references is among the paths deleted on the branch, given to it as input (the paths deleted in merge-base..HEAD). A capability none of whose missing paths the branch deleted is not reported. It reads only the model, an fs.FS and the deleted-path list, runs no git, and changes no file outside internal/modellint.

## Links

### Out

- depends-on [modellint-stale-refs](rhiza/workflows/built-in-workflows/built-in-review/modellint-stale-refs.md)
- decided-by [lint-hits-feed-model-check](../decision/lint-hits-feed-model-check.md)

### In

- depends-on [code-check-runs-lints](rhiza/workflows/built-in-workflows/built-in-review/code-check-runs-lints.md)
- verifies [modellint-branch-removed-code-verified](../evidence/modellint-branch-removed-code-verified.md)
- verifies [modellint-branch-removed-code-qa-1](../evidence/modellint-branch-removed-code-qa-1.md)
- verifies [pr-569-model-lints-in-review](../evidence/pr-569-model-lints-in-review.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
