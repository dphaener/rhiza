# milestone-number-reserved-qa-1

evidence · `01M77BEJ4503JSZVZ5VN28YNS4` · actor [claude-code-215c18e1](../actor/claude-code-215c18e1.md)

## result

PASS

## statement

QA claude-code-215c18e1 at 8e55f6d4a, by hand with the built binary (go build -o scratchpad/rhiza-qa ./cmd/rhiza) in a scratch Rhiza root and repo with a stub plan-milestone workflow. rhiza run start --foreground --workflow plan --new X --milestone --json reserved 1, 2, 3 (Ship v5, which failed) and then 4, so the failed run's number was skipped. After landing milestone-12 and milestone-50 with --record, the next starts drew 51, 52, 53 and 54. Four run starts launched at once in the background drew 55, 56, 57 and 58, all different. rhiza run status <id> --json shows milestone_number for each. A plan run on an intention (--new with no --milestone) has no milestone_number. Scoped tests pass: go test -trimpath ./cmd/rhiza -run 'TestScripts/(plan-milestone-number|plan-run-milestone-env|plan-milestone-check|describe-golden)$', ./internal/run, and ./internal/api -run Milestone.

## Links

### Out

- verifies [milestone-number-reserved](../capability/milestone-number-reserved.md)

## History

- `01M7774MAPJ7VCPN1Q9M3RHMV4` Fix Milestone plans running at the same time pick the same milestone number and handle, so the second to merge fails CI self-check with a duplicate handle. Evidence: the agent docs plan (run 01M76EC9, record 01M76ECB) named its intent milestone-32, already held by the cleanup milestone 01M7640B, and failed pr-green-infra on PR #456. I renamed it to milestone-34 by hand (record cs 01M76V0W). The efficient build cycle plan (run 01M76STX, record 01M76SV0) had meanwhile also chosen milestone-34 and failed CI on PR #466 ('duplicate handle: node 01M76NWA... has handle milestone-34, already held by node 01M76XMNA9...'); renamed to milestone-35 by hand (record cs 01M774VK) and resume --retry. Cause: the planner takes the next number from its own branch's model, which can't see other in-flight milestone branches. Fix ideas: the daemon assigns the milestone number when it opens the record change set (counting open milestone records too), or the plan's close-out/rebase renumbers on conflict. The pre-push self-check (01M76R1E) would catch it locally only after merging main.
