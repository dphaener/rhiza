# modellint-stale-refs

capability · `01M7AZSB7XH20CA9CJCEATG154` · parent [built-in-review](../built-in-review.md)

## acceptance

go test -trimpath ./internal/modellint/... passes with table tests in internal/modellint/refs_test.go: a leaf naming a present path, Go test and testscript reports nothing; a missing path, a missing Go test and a missing testscript each report one hit naming the node, the field and the reference; a capability whose every path is gone is reported as unretired; failure cases: a node the change set did not write, naming a missing path, reports nothing, and a token whose first segment is no top-level entry ("and/or", a URL) is not a reference.

## statement

A new pure package internal/modellint reports, for each node the change set adds or changes in its to-be view (its leaves and the other capabilities and nodes it writes), each reference in its statement or acceptance that the head tree lacks: a repo path (a token whose first segment, after stripping a leading `./` and a trailing `/...`, is a top-level entry of the head tree), a Go test name `TestXxx` with no `func TestXxx(` in any *_test.go, and a testscript `TestScripts/<name>` with no `<name>.txtar` in the tree; a capability whose path references are all gone is reported as an unretired capability. It reads only the model and an fs.FS, runs no git, and does not touch internal/check or any file outside internal/modellint.

## Links

### Out

- decided-by [lint-hits-feed-model-check](../../../../../decision/lint-hits-feed-model-check.md)

### In

- depends-on [modellint-unapplied-decisions](modellint-unapplied-decisions.md)
- depends-on [modellint-branch-removed-code](../../../../modellint-branch-removed-code.md)
- verifies [modellint-stale-refs-verified](../../../../../evidence/modellint-stale-refs-verified.md)
- verifies [modellint-stale-refs-qa-1](../../../../../evidence/modellint-stale-refs-qa-1.md)
- verifies [pr-569-model-lints-in-review](../../../../../evidence/pr-569-model-lints-in-review.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
