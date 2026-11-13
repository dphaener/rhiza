# plan-merge-loop-cap-six

decision · `01M7AQEYHR99BXC70GVCCW2Q06`

## reasoning

Each rebased head waits 2-7 minutes for CI, and on a main where several plan runs merge render PRs, another merge often lands in that window. Raising the cap to ship-change's keeps the two workflows' merge loops alike and follows the earlier Decision's rule; serialising merges would be a new mechanism outside this bug's scope.

## statement

The plan and plan-milestone workflows cap their merge loop's steps (rebase, wait-rebased, rebased-green and merge) at 6 visits each, ship-change's current cap for the same loop, so a plan PR on a busy main gets six rebase-and-merge tries before the run ends failed; the other steps' caps are unchanged

## wording

(derived) plan-merge-loops-to-rebase gave the plan workflows ship-change's cap; ship-change's merge loop is now 6 visits and the plan workflows kept 3, which run 01M7AH6RHKMZRAHGJKJFM4VTNT exhausted on a busy main

## Links

### Out

- supersedes [plan-merge-loops-to-rebase](plan-merge-loops-to-rebase.md)
- derived-from [plan-merge-race-research](../source/plan-merge-race-research.md)

### In

- decided-by plan-run-merge-rebases (retired)
- decided-by [plan-merge-loop-six-tries](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-merge-loop-six-tries.md)

## History

- `01M7AQ89397MPNWGFSAYG3BK77` Fix watchdog unmerged PR at step rebase: run 01M7AH6RHKMZRAHGJKJFM4VTNT; log: PR green at 6ce8f26 and unmerged after step rebase failed: step "rebase" went over max_visits 3
