# question-b8f88971-1

question · `01M73G61GENMB6SV7YM0PCCZ56`

## session

b8f88971-6f9a-4c0c-b545-01a1c700dbf9

## statement

Run 01M73CQ4 looped merge -> gate-rebase-merge until max_visits because PR #426 had been converted to draft after ship opened it (gh: 'Pull Request is still a draft'), while ship check still called it green. When rhiza ship merge finds the run's PR is a draft, what should it do?

## Links

### In

- answers [decision-draft-pr-is-a-hold](../decision/decision-draft-pr-is-a-hold.md)

## History

- `01M73G0WSCCBG7H96NDY1H4KQ6` Fix watchdog unmerged PR at step gate-rebase-merge: run 01M73CQ4D71CW038ZBNSQB1SJF; log: PR green at 0b28569 and unmerged after step gate-rebase-merge failed: step "gate-rebase-merge" went over max_visits 6
