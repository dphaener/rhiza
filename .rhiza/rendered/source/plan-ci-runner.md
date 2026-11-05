# plan-ci-runner

source · `01M6HWK1AWW9CKCY5PB01YDASY`

## quotes

```json
[
  {
    "by": "owner",
    "text": "In the past this was because the remote origin was https and not ssh.",
    "where": "relayed by the orchestrator in a run message during the interview (2026-11-03T23:47:10Z), answering question 4 (how the workflow edit lands)"
  }
]
```

## raw_sha256

0e5cfab25c835127a9e23b2cabe41b704d73ffc84907231abfbd87dbb7c093a6

## session

226cf21d-db5d-4a1f-82be-01a1a3c4134c

## summary

Autonomous plan run for change set 'Update github runner to'. No owner was present; the planner raised one question (serves, main CI, ci install template, workflow push access) and the arbiter answered with Decision ci-main-and-cancel. Settled: the change set serves milestone-25; ship-change's wait-main stops waiting on GitHub and instead compares the merge commit's tree with the tested PR head's tree locally, failing with a stop: line on a mismatch; the rhiza ci install template gains a per-PR concurrency group with cancel-in-progress but keeps its push-to-main check. The arbiter first ruled (part 4 of ci-main-and-cancel) that agent leaves touch only Go code and the .github/workflows/ci.yml edit is left to the owner to push by hand, because the agents' gh token lacks the workflow scope. During the interview the orchestrator relayed that push access is fixed: origin now uses SSH with the VM's key, and the owner said "In the past this was because the remote origin was https and not ssh." Decision ci-yml-agent-built supersedes part 4: the ci.yml edit is an agent-built leaf of this change set.

## title

Plan: cancel stale CI runs, stop re-testing on main

## Links

### In

- derived-from [ci-main-and-cancel](../decision/ci-main-and-cancel.md)
- derived-from [ship-main-tree-check](../capability/rhiza/workflows/workflow-runs/ship-main-tree-check.md)
- derived-from [ci-install-cancels-stale](../capability/rhiza/integrity/ci-install-cancels-stale.md)
- derived-from [ci-yml-agent-built](../decision/ci-yml-agent-built.md)
- derived-from [ci-yml-cancels-stale](../capability/rhiza/integrity/ci-yml-cancels-stale.md)
- derived-from [ci-main-landing-order](../decision/ci-main-landing-order.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
