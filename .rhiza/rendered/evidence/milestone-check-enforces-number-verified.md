# milestone-check-enforces-number-verified

evidence · `01M77B1PHV3AKC6YZS89107BYY` · actor [claude-code-ab5f8f28](../actor/claude-code-ab5f8f28.md)

## result

PASS

## statement

Verifier at d05f25d34: go test -trimpath ./internal/api/ -run TestRunMilestoneCheckEnforcesNumber -count=1 -v passes. The test starts plan-milestone runs with MilestoneNumber 7 and calls run.milestone-check through the API registry. It passes for milestone-7 with field milestone=7 and done point m7-first. It fails for handle milestone-8, for field milestone=8, and for unprefixed done point "first". Each error names the intent, contains "milestone number 7" and has a Fix line that names milestone-7, milestone=7 and the m7- prefix. The full internal/api package run printed PASS for every test, then the root-state guard reported .rhiza-root/usage.live.json changed. That file is written by the live root, not by this diff. go vet and gofmt are clean. In the diff (newMilestoneNumberProblems), only new intents that carry a milestone field (absent from the landed disk model) are checked, and only when the run reserved a number.

## Links

### Out

- verifies [milestone-check-enforces-number](../capability/milestone-check-enforces-number.md)

## History

- `01M7774MAPJ7VCPN1Q9M3RHMV4` Fix Milestone plans running at the same time pick the same milestone number and handle, so the second to merge fails CI self-check with a duplicate handle. Evidence: the agent docs plan (run 01M76EC9, record 01M76ECB) named its intent milestone-32, already held by the cleanup milestone 01M7640B, and failed pr-green-infra on PR #456. I renamed it to milestone-34 by hand (record cs 01M76V0W). The efficient build cycle plan (run 01M76STX, record 01M76SV0) had meanwhile also chosen milestone-34 and failed CI on PR #466 ('duplicate handle: node 01M76NWA... has handle milestone-34, already held by node 01M76XMNA9...'); renamed to milestone-35 by hand (record cs 01M774VK) and resume --retry. Cause: the planner takes the next number from its own branch's model, which can't see other in-flight milestone branches. Fix ideas: the daemon assigns the milestone number when it opens the record change set (counting open milestone records too), or the plan's close-out/rebase renumbers on conflict. The pre-push self-check (01M76R1E) would catch it locally only after merging main.
