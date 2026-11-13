# idle-poll-budget-judges-no-evidence

decision · `01M7B8BJ5V9P5G0JQRGYEGBDD4`

## reasoning

CI's budgets timing log shows the idle-poll median stepped from about 1.2-2.5 ms to 2.4-4.8 ms (6.0 ms in PR #567's failing attempt, 3% CPU pressure) exactly at PR #516's merge, while every read budget stayed far under 50 ms; the added cost is the forced catch-up the fixture's per-poll Evidence call now makes, which an idle driver never pays. Raising the budget would hide a real idle regression later; quarantining would drop TestReadBudget's whole read budget; weakening settle() would reopen the headless-run-agent-tabs race #516 fixed.

## statement

Fix the flake in the measure, not the product or the budget: driverPollBudget's idle measure times a poll without a verifier's Evidence judgment, because since PR #516 every Evidence read forces an index catch-up (settle) and the real driver judges Evidence only once, when a verifier session ends (judgeVerifier), never on an idle poll. idlePollCPUBudget stays 5 ms, unrelatedPollCPUBudget 50 ms, and #516's forced catch-up in stepRunner.Evidence stays as it is

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [readbudget-driver-poll-flake-plan](../source/readbudget-driver-poll-flake-plan.md)

### In

- decided-by [idle-poll-budget-times-idle-poll](../capability/rhiza/orchestration/daemon/idle-poll-budget-times-idle-poll.md)

## History

- `01M7B7P8D389GB7N94YPTXAQCS` Fix flaky test TestReadBudget: failed then passed on a rerun in PR #567 at head c5cc453c8b33b60fd6078c0aeaf5e1331ce3239f; reproduce it locally, repeatedly, then fix it or quarantine it
