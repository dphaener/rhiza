# planner-gets-milestone-number-verified

evidence · `01M779YPA76M9D7MS9GVQWNH8G` · actor [claude-code-04fca9b2](../actor/claude-code-04fca9b2.md)

## result

PASS

## statement

Verifier claude-code-04fca9b2 at 1e6517862: read the diff (session.Record.MilestoneNumber, SessionStartRequest.milestoneNumber set by stepRunner.StartAgent from the run state's reserved MilestoneNumber, launchEnv sets RHIZA_MILESTONE_NUMBER=N when >0 else unsets it; plan.md Milestone section tells the planner to name the intent milestone-N, field milestone=N, done points prefixed mN-, never counting milestones from its own branch). Ran go test -trimpath ./internal/api -run 'TestLaunchEnvMilestoneNumber|TestSessionRun|MilestoneNumber' -count=1: ok (TestLaunchEnvMilestoneNumber checks set with MilestoneNumber 36 and unset with none). Ran go test -trimpath ./internal/guide/... ./internal/session/... -count=1: ok (TestPlanMilestoneNumber checks RHIZA_MILESTONE_NUMBER, milestone-N, milestone=N, `mN-`, Never count milestones). Full ./internal/api with -timeout 25m: every test PASS; the package-level guard then flagged ~/code/.rhiza-root/usage.live.json changed, which is the live root daemon's usage file, untouched by this diff.

## Links

### Out

- verifies [planner-gets-milestone-number](../capability/planner-gets-milestone-number.md)

## History

- `01M7774MAPJ7VCPN1Q9M3RHMV4` Fix Milestone plans running at the same time pick the same milestone number and handle, so the second to merge fails CI self-check with a duplicate handle. Evidence: the agent docs plan (run 01M76EC9, record 01M76ECB) named its intent milestone-32, already held by the cleanup milestone 01M7640B, and failed pr-green-infra on PR #456. I renamed it to milestone-34 by hand (record cs 01M76V0W). The efficient build cycle plan (run 01M76STX, record 01M76SV0) had meanwhile also chosen milestone-34 and failed CI on PR #466 ('duplicate handle: node 01M76NWA... has handle milestone-34, already held by node 01M76XMNA9...'); renamed to milestone-35 by hand (record cs 01M774VK) and resume --retry. Cause: the planner takes the next number from its own branch's model, which can't see other in-flight milestone branches. Fix ideas: the daemon assigns the milestone number when it opens the record change set (counting open milestone records too), or the plan's close-out/rebase renumbers on conflict. The pre-push self-check (01M76R1E) would catch it locally only after merging main.
