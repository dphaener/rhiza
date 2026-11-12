# milestone-check-enforces-number-verified-2

evidence · `01M77CHH07ZE57CYX44BB2GMAZ` · actor [claude-code-126f8a33](../actor/claude-code-126f8a33.md)

## result

PASS

## statement

Verified at b2455268f (attempt 2). Read the diff: newMilestoneNumberProblems now treats as new the milestone intents the plan's change sets created (landed or not) and runs over every plan change set, so a record landed during the interview is still checked, which addresses the confirmed quality/model finding. In a scratch worktree outside ~/code (inside it the root daemon's usage.live.json trips rootguard, unrelated to this change), go test -trimpath ./internal/api -run Milestone passes, including TestRunMilestoneCheckEnforcesNumber: reserved number 7 used passes; wrong handle milestone-8, wrong field milestone=8 and unprefixed done point 'first' each fail naming the intent, 'milestone number 7' and a Fix: line; a misnumbered milestone-8 still fails after the record is proposed, accepted and landed. With run_milestone.go reverted to HEAD~1 that new case fails ('landed record: passed; want a failure'), so the test pins the fix. TestScripts/plan-milestone-number and plan-milestone-check also pass.

## Links

### Out

- verifies [milestone-check-enforces-number](../capability/milestone-check-enforces-number.md)

## History

- `01M7774MAPJ7VCPN1Q9M3RHMV4` Fix Milestone plans running at the same time pick the same milestone number and handle, so the second to merge fails CI self-check with a duplicate handle. Evidence: the agent docs plan (run 01M76EC9, record 01M76ECB) named its intent milestone-32, already held by the cleanup milestone 01M7640B, and failed pr-green-infra on PR #456. I renamed it to milestone-34 by hand (record cs 01M76V0W). The efficient build cycle plan (run 01M76STX, record 01M76SV0) had meanwhile also chosen milestone-34 and failed CI on PR #466 ('duplicate handle: node 01M76NWA... has handle milestone-34, already held by node 01M76XMNA9...'); renamed to milestone-35 by hand (record cs 01M774VK) and resume --retry. Cause: the planner takes the next number from its own branch's model, which can't see other in-flight milestone branches. Fix ideas: the daemon assigns the milestone number when it opens the record change set (counting open milestone records too), or the plan's close-out/rebase renumbers on conflict. The pre-push self-check (01M76R1E) would catch it locally only after merging main.
