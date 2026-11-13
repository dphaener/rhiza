# readbudget-driver-poll-flake-plan

source · `01M7B8AM8RXZ0MVMXEG656Q3HB`

## quotes

```json
[]
```

## raw_sha256

8a526f4c1ef11926aa1da3d4889007c3003975210915b6a7bc05dd87eba60905

## session

173c935e-3487-4963-ab92-01a1d67beda3

## summary

Autonomous plan run, no owner. CI run 43495062180 attempt 1 (PR #567, head c5cc453c8) failed the budgets job on one line of TestReadBudget: 'driver poll idle: median 6.004ms of CPU over 5ms' (driverPollBudget in internal/api/run_model_cache_test.go, idlePollCPUBudget = 5 ms) at 3% CPU pressure, with every read budget well under 50 ms; attempt 2 passed at 2.137 ms. The planner read the idle-poll medians from the budgets job's timing log of recent CI runs: about 1.2-2.5 ms before 2026-11-13T01:27Z, about 2.4-4.8 ms after (6.0 ms once), a step at the merge of PR #516 (d2749ed50, Decision verifier-evidence-read-fresh), which made stepRunner.Evidence call settle() so every Evidence read forces an index catch-up (Refresh) past the poll debounce. The budget's fixture poll (driverFixture.poll) judges a verifier's Evidence on every poll, but the real driver calls Evidence only in judgeVerifier, once when a verifier session ends (internal/run/leaves.go), never on an idle poll. So the idle budget now times a forced catch-up an idle driver never makes, leaving almost no headroom under 5 ms; the product's idle poll did not regress. The planner chose to fix the measure (time the idle poll without the Evidence judgment), keeping both budgets and #516's forced catch-up, over raising the budget, quarantining, or weakening settle(). Not breaking, not severe.

## title

Plan of the TestReadBudget driver-poll flake (autonomous)

## Links

### In

- derived-from [idle-poll-budget-judges-no-evidence](../decision/idle-poll-budget-judges-no-evidence.md)

## History

- `01M7B7P8D389GB7N94YPTXAQCS` Fix flaky test TestReadBudget: failed then passed on a rerun in PR #567 at head c5cc453c8b33b60fd6078c0aeaf5e1331ce3239f; reproduce it locally, repeatedly, then fix it or quarantine it
