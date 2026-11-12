# planner-gets-milestone-number

capability · `01M777B9HFVF828QJGDWM3M4CH`

## acceptance

A test in internal/api checks the planner session's environment carries RHIZA_MILESTONE_NUMBER for a new-milestone plan run and omits it otherwise; a test in internal/guide checks the plan guide's Milestone section names RHIZA_MILESTONE_NUMBER, milestone-N, milestone=N and the mN- prefix

## statement

The planner of a milestone plan run gets the reserved number as RHIZA_MILESTONE_NUMBER, and the plan guide tells it to name the milestone intent milestone-N with field milestone=N and its done points prefixed mN-, never counting milestones from its own branch

## Links

### Out

- decided-by [milestone-number-reserved-at-run-start](../decision/milestone-number-reserved-at-run-start.md)

### In

- verifies [planner-gets-milestone-number-verified](../evidence/planner-gets-milestone-number-verified.md)
- verifies [planner-gets-milestone-number-qa-1](../evidence/planner-gets-milestone-number-qa-1.md)
- verifies [pr-474-fix-milestone-plans-running-at](../evidence/pr-474-fix-milestone-plans-running-at.md)

## History

- `01M7774MAPJ7VCPN1Q9M3RHMV4` Fix Milestone plans running at the same time pick the same milestone number and handle, so the second to merge fails CI self-check with a duplicate handle. Evidence: the agent docs plan (run 01M76EC9, record 01M76ECB) named its intent milestone-32, already held by the cleanup milestone 01M7640B, and failed pr-green-infra on PR #456. I renamed it to milestone-34 by hand (record cs 01M76V0W). The efficient build cycle plan (run 01M76STX, record 01M76SV0) had meanwhile also chosen milestone-34 and failed CI on PR #466 ('duplicate handle: node 01M76NWA... has handle milestone-34, already held by node 01M76XMNA9...'); renamed to milestone-35 by hand (record cs 01M774VK) and resume --retry. Cause: the planner takes the next number from its own branch's model, which can't see other in-flight milestone branches. Fix ideas: the daemon assigns the milestone number when it opens the record change set (counting open milestone records too), or the plan's close-out/rebase renumbers on conflict. The pre-push self-check (01M76R1E) would catch it locally only after merging main.
