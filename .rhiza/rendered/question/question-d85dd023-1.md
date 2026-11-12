# question-d85dd023-1

question · `01M77788WM4669KTMNJDY0T7B3`

## session

d85dd023-7452-4f12-87c3-01a1ce728739

## statement

Which intent should bug change set 01M7774 (concurrent milestone plans pick the same milestone number and handle) serve? The earlier bug fix 'A queued milestone request waits for another milestone' served milestone-25 (the backlog is empty: the daemon builds and ships the queued work on its own). / How should concurrent milestone plans stop picking the same number and handle? / Should run.milestone-check (plan step before review) also fail a plan whose milestone intent doesn't use the reserved number (handle milestone-N, field milestone=N, done points prefixed mN-)?

## Links

### In

- answers [milestone-number-reserved-at-run-start](../decision/milestone-number-reserved-at-run-start.md)

## History

- `01M7774MAPJ7VCPN1Q9M3RHMV4` Fix Milestone plans running at the same time pick the same milestone number and handle, so the second to merge fails CI self-check with a duplicate handle. Evidence: the agent docs plan (run 01M76EC9, record 01M76ECB) named its intent milestone-32, already held by the cleanup milestone 01M7640B, and failed pr-green-infra on PR #456. I renamed it to milestone-34 by hand (record cs 01M76V0W). The efficient build cycle plan (run 01M76STX, record 01M76SV0) had meanwhile also chosen milestone-34 and failed CI on PR #466 ('duplicate handle: node 01M76NWA... has handle milestone-34, already held by node 01M76XMNA9...'); renamed to milestone-35 by hand (record cs 01M774VK) and resume --retry. Cause: the planner takes the next number from its own branch's model, which can't see other in-flight milestone branches. Fix ideas: the daemon assigns the milestone number when it opens the record change set (counting open milestone records too), or the plan's close-out/rebase renumbers on conflict. The pre-push self-check (01M76R1E) would catch it locally only after merging main.
