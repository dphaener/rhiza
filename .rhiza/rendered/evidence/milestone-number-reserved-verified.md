# milestone-number-reserved-verified

evidence · `01M77A2FPVEA6Q323GYARBW9PJ` · actor [claude-code-546189ed](../actor/claude-code-546189ed.md)

## result

PASS

## statement

Verified at 8eb70a774 (attempt 2). Read the diff of 4ddef2aa9 and 8eb70a774: runRunStart reserves a number only when newMilestone (milestone plan with no request milestone, or a plan on a record change set), drawing it under a lock as one more than the highest landed milestone field and every run journal's MilestoneNumber (ended and failed included), and journals it in Started.MilestoneNumber, which run status shows as milestone_number. Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(plan-milestone-number|describe-golden)$': PASS. plan-milestone-number starts new-milestone plan runs with run start against landed milestones 12 and 35 and sees 36 then 37, a failed run holding 38 so the next gets 39, run status --json showing milestone_number, and a milestone request naming milestone-35 started with no milestone_number. Ran go test -trimpath ./internal/run -run TestFoldMilestoneNumber (failed run folds 36, no number folds 0 and omits the field): PASS. Ran go test -trimpath ./internal/api -run Milestone: all 17 tests PASS, including TestReserveMilestoneNumber and TestLaunchEnvMilestoneNumber. That run then hit the enclosing-root state-folder guard. Rerunning the leaf's api tests alone passed cleanly, so the guard failure does not come from this change. Attempt 1's three findings (no run-start test, no existing-milestone test, no internal/run test) are each addressed.

## Links

### Out

- verifies [milestone-number-reserved](../capability/milestone-number-reserved.md)

## History

- `01M7774MAPJ7VCPN1Q9M3RHMV4` Fix Milestone plans running at the same time pick the same milestone number and handle, so the second to merge fails CI self-check with a duplicate handle. Evidence: the agent docs plan (run 01M76EC9, record 01M76ECB) named its intent milestone-32, already held by the cleanup milestone 01M7640B, and failed pr-green-infra on PR #456. I renamed it to milestone-34 by hand (record cs 01M76V0W). The efficient build cycle plan (run 01M76STX, record 01M76SV0) had meanwhile also chosen milestone-34 and failed CI on PR #466 ('duplicate handle: node 01M76NWA... has handle milestone-34, already held by node 01M76XMNA9...'); renamed to milestone-35 by hand (record cs 01M774VK) and resume --retry. Cause: the planner takes the next number from its own branch's model, which can't see other in-flight milestone branches. Fix ideas: the daemon assigns the milestone number when it opens the record change set (counting open milestone records too), or the plan's close-out/rebase renumbers on conflict. The pre-push self-check (01M76R1E) would catch it locally only after merging main.
