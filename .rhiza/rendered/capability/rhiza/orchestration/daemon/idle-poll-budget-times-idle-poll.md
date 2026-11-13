# idle-poll-budget-times-idle-poll

capability · `01M7B8BJ8E3YCTNWCHPT45HN4P` · parent [daemon](../daemon.md)

## acceptance

Reproduce first: go test -trimpath -count=1 -p 1 -run '^TestReadBudget$' ./internal/api with RHIZA_BUDGET_LOG set to a scratch file, run 3 times one after another (never in parallel and never under generated CPU load, per AGENTS.md), recording each run's 'driver poll idle' median and CPU pressure in the Checks: block. Then in driverPollBudget the idle measure polls without judging Evidence (f.poll with no leaf), the warm-up polls and the 'after an unrelated entry' measure stay as they are, idlePollCPUBudget stays 5 ms and unrelatedPollCPUBudget 50 ms, the comment above the idle measure says why an idle poll judges no Evidence, and stepRunner.Evidence and settle() are unchanged. The same 3 runs then pass, each idle median recorded, their middle value lower than before on the same machine; go test -trimpath -count=1 -run '^(TestDriverPollIsCheap|TestEvidenceReadsPastSessionWrites)$' ./internal/api passes. Only internal/api/run_model_cache_test.go changes.

## statement

The driver CPU budget in TestReadBudget (driverPollBudget, internal/api/run_model_cache_test.go) times an idle poll as a real idle driver makes it, the leaves, answers and session status with no verifier Evidence judgment, so the forced catch-up Evidence makes when a verifier ends no longer pushes the idle median against its 5 ms budget

## Links

### Out

- decided-by [idle-poll-budget-judges-no-evidence](../../../../decision/idle-poll-budget-judges-no-evidence.md)

### In

- verifies [idle-poll-budget-times-idle-poll-verified](../../../../evidence/idle-poll-budget-times-idle-poll-verified.md)
- verifies [idle-poll-budget-times-idle-poll-qa-1](../../../../evidence/idle-poll-budget-times-idle-poll-qa-1.md)
- verifies [pr-570-fix-flaky-test-testreadbudget](../../../../evidence/pr-570-fix-flaky-test-testreadbudget.md)

## History

- `01M7B7P8D389GB7N94YPTXAQCS` Fix flaky test TestReadBudget: failed then passed on a rerun in PR #567 at head c5cc453c8b33b60fd6078c0aeaf5e1331ce3239f; reproduce it locally, repeatedly, then fix it or quarantine it
