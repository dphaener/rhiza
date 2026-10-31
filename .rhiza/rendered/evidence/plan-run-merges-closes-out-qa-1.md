# plan-run-merges-closes-out-qa-1

evidence · `01M695829QE9E73MBRD0V2KP2Q` · actor [claude-code-418647f6](../actor/claude-code-418647f6.md)

## result

PASS

## statement

QA cycle 1 at a3dd859: built rhiza with go build to a scratch path and ran it by hand in a scratch repo with a bare origin.git, a fake rex/claude (roles doing real rhiza writes) and a fake gh (no network). With green checks the plan run ran wait-pr, pr-green, merge ('pr merge 7 --squash --match-head-commit <pushed head>', no --admin/--auto/--delete-branch) and close-out: worktree repo.plan-reports and cs/plan-reports gone locally and on origin, main fast-forwarded to the merged head, .git/rhiza/ship empty; the planned change set stayed proposed (run landed nothing). A following ship-change run on that change set opened PR #8 and its merge step called 'pr merge 8 --squash --match-head-commit <head>', reported a merge_commit and no already_merged; ship state then held pr 8.

## Links

### Out

- verifies [plan-run-merges-closes-out](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-merges-closes-out.md)

## History

- `01M625RCFY4RB9HC0B69Y08VFV` A plan run ends with its plan on a pushed branch and an open PR, as ship-change does, so no one renders, commits or opens the plan's PR by hand.
