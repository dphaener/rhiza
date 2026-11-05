# ci-yml-agent-built

decision · `01M6HX39EBGJBJJS89P3N86CGX`

## quote

In the past this was because the remote origin was https and not ssh.

## statement

The .github/workflows/ci.yml edit is an agent-built leaf of this change set, replacing part (4) of ci-main-and-cancel: origin now uses SSH with the VM's key, which the token's missing workflow scope does not limit. It ships in the same PR as ship-main-tree-check; Decision ci-main-landing-order says how this change set's own ship gets past wait-main, which the installed binary still runs on CI for main until the daemon's self-update replaces it.

## Links

### Out

- supersedes [ci-main-and-cancel](ci-main-and-cancel.md) — Supersedes part (4) only: the ci.yml edit is an agent-built leaf, not a hand push by the owner. Parts (1) to (3) stand.
- derived-from [plan-ci-runner](../source/plan-ci-runner.md)

### In

- decided-by [ci-yml-cancels-stale](../capability/rhiza/integrity/ci-yml-cancels-stale.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
