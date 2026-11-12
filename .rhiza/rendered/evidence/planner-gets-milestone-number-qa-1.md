# planner-gets-milestone-number-qa-1

evidence · `01M77BEJH63J749BR01WDH2RQ9` · actor [claude-code-215c18e1](../actor/claude-code-215c18e1.md)

## result

PASS

## statement

QA claude-code-215c18e1 at 8e55f6d4a, with the built binary, a fake rex that runs the agent command and a fake claude that writes its RHIZA_MILESTONE* environment to a file. A plan-milestone run whose agent interview step reserved 61 launched the planner with RHIZA_MILESTONE_NUMBER=61. A plan run on an intention launched its planner with no RHIZA_MILESTONE_NUMBER. rhiza guide plan (built binary) prints the Milestone bullet: name the intent from RHIZA_MILESTONE_NUMBER (N), handle milestone-N, field milestone=N, done points prefixed mN-, never counting milestones from your own branch. go test -trimpath ./internal/guide ./internal/session and ./internal/api -run Milestone pass.

## Links

### Out

- verifies [planner-gets-milestone-number](../capability/planner-gets-milestone-number.md)

## History

- `01M7774MAPJ7VCPN1Q9M3RHMV4` Fix Milestone plans running at the same time pick the same milestone number and handle, so the second to merge fails CI self-check with a duplicate handle. Evidence: the agent docs plan (run 01M76EC9, record 01M76ECB) named its intent milestone-32, already held by the cleanup milestone 01M7640B, and failed pr-green-infra on PR #456. I renamed it to milestone-34 by hand (record cs 01M76V0W). The efficient build cycle plan (run 01M76STX, record 01M76SV0) had meanwhile also chosen milestone-34 and failed CI on PR #466 ('duplicate handle: node 01M76NWA... has handle milestone-34, already held by node 01M76XMNA9...'); renamed to milestone-35 by hand (record cs 01M774VK) and resume --retry. Cause: the planner takes the next number from its own branch's model, which can't see other in-flight milestone branches. Fix ideas: the daemon assigns the milestone number when it opens the record change set (counting open milestone records too), or the plan's close-out/rebase renumbers on conflict. The pre-push self-check (01M76R1E) would catch it locally only after merging main.
