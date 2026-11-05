# question-226cf21d-1

question · `01M6HWCD54ZEEPYXMAZN20DVZK`

## session

226cf21d-db5d-4a1f-82be-01a1a3c4134c

## statement

Change set 'Update github runner to' (cancel stale CI runs, stop re-testing on main) serves no intent yet. Which intent should it serve? / Stopping CI on push to main breaks ship-change: its wait-main step (`rhiza ship check --main --settled`) waits for runs on the merge commit and would time out after 1h on every ship. Since ship-change rebases a PR main has moved past and merges only the head whose checks passed, the merge commit's tree was already tested. How should main be handled? / Should `rhiza ci install`'s workflow template (internal/api/ci.go, for other repos) change too, or only this repo's .github/workflows/ci.yml? / The gh token the agents push with lacks the `workflow` scope, so GitHub rejects any push touching .github/workflows/ci.yml. How does the workflow edit land?

## Links

### In

- answers [ci-main-and-cancel](../decision/ci-main-and-cancel.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
