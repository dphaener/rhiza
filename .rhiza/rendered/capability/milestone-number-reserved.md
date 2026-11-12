# milestone-number-reserved

capability · `01M777B9G3TT1SDCSTGG2Z820G`

## acceptance

Tests in internal/api (and internal/run for the journal field) start two milestone plan runs back to back with landed milestones up to 35 and see 36 then 37; a reservation held only by an ended or failed run is still skipped; a run planning an existing milestone gets no number; rhiza run status --json shows the reserved number

## statement

Starting a plan run for a new milestone (a milestone request, or run start --new --milestone) reserves its milestone number N, one more than the highest of every landed milestone intent's milestone field and every number reserved in any run journal in the shared run store (ended and failed runs included), stores N in the run's journal and shows it in rhiza run status; a run that plans an existing milestone reserves nothing

## Links

### Out

- decided-by [milestone-number-reserved-at-run-start](../decision/milestone-number-reserved-at-run-start.md)

### In

- verifies [milestone-number-reserved-verified](../evidence/milestone-number-reserved-verified.md)
- verifies [milestone-number-reserved-qa-1](../evidence/milestone-number-reserved-qa-1.md)
- verifies [pr-474-fix-milestone-plans-running-at](../evidence/pr-474-fix-milestone-plans-running-at.md)

## History

- `01M7774MAPJ7VCPN1Q9M3RHMV4` Fix Milestone plans running at the same time pick the same milestone number and handle, so the second to merge fails CI self-check with a duplicate handle. Evidence: the agent docs plan (run 01M76EC9, record 01M76ECB) named its intent milestone-32, already held by the cleanup milestone 01M7640B, and failed pr-green-infra on PR #456. I renamed it to milestone-34 by hand (record cs 01M76V0W). The efficient build cycle plan (run 01M76STX, record 01M76SV0) had meanwhile also chosen milestone-34 and failed CI on PR #466 ('duplicate handle: node 01M76NWA... has handle milestone-34, already held by node 01M76XMNA9...'); renamed to milestone-35 by hand (record cs 01M774VK) and resume --retry. Cause: the planner takes the next number from its own branch's model, which can't see other in-flight milestone branches. Fix ideas: the daemon assigns the milestone number when it opens the record change set (counting open milestone records too), or the plan's close-out/rebase renumbers on conflict. The pre-push self-check (01M76R1E) would catch it locally only after merging main.
