# rebase-when-behind

capability · `01M5X3GCHZRJ9YPV8BA37J1AN3` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript with a local bare origin where main moves during the run: once with a rendered-only conflict (the run rebases, re-renders and merges), once with a code conflict (the step fails and a question names the files).

## statement

Before landing and again before merging, ship-change rebases a PR that main has moved past: conflicts only under .rhiza/rendered/ are settled by re-rendering; any other conflict fails the step and becomes a question; after a rebase it pushes and waits for CI again.

## Links

### Out

- decided-by [ship-change-holds-up](../../../../decision/ship-change-holds-up.md)

### In

- verifies [rebase-when-behind-tests](../../../../evidence/rebase-when-behind-tests.md)

## History

- `01M5X3GBH4ZQTZVXFEHR966747` ship-change that holds up: the agent step passes only when the agent says it is done, land renders with the branch's own build in Rhiza's repo, a PR that main has moved past is rebased, and a failed run can be resumed from the step that failed.
