# t5-d13-budget-threshold

decision · `01M5SWS2KCZEMS8K7MR95WPJCS`

## accepted_by

owner

## revised

revised by the orchestrator under the owner's authorization: the budget test's limit is 50 ms, milestone 1's read target, matching the read and slice budgets. TestBudget failed on main after PR #10 (fresh-read median 27.4 ms on GitHub's shared runner against the 25 ms limit, after passing on the PR); a limit that close to the runner's normal speed fails at random. Fixed in PR #11 (3b5d1da); the benchmarks still report the tighter numbers.

## settled_by

Agent's option, accepted by the owner: Median < 25 ms (Recommended): Half the 50 ms target over ~20 runs of the fresh-read path (open, freshness check, a Resolve + Ancestors + Children + Links); room to spare so CI doesn't flake, still fails a real regression. Skipped under -short.

## statement

Threshold: median < 50 ms over 20 runs of the fresh-read path (milestone 1's read target; revised from 25 ms in PR #11); skipped under `-short`.

## wording

agent of milestone-1 task 5 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-5](../source/plan-task-5.md)
- derived-from [task-interviews](../source/task-interviews.md)
- derived-from [founding-conversation](../source/founding-conversation.md)

### In

- decided-by [fast-reads](../capability/rhiza/views/fast-reads.md)

## History

- `01M5SWS22WXZMST7XXR9ZKY4EE` Task 5's plan decisions (PR #7, merge ed8d292), accepted by the owner
