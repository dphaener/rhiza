# milestone-check-enforces-number-qa-1

evidence · `01M77BEJS0T5FXSGYH84CY6MPP` · actor [claude-code-215c18e1](../actor/claude-code-215c18e1.md)

## result

PASS

## statement

QA claude-code-215c18e1 at 8e55f6d4a, with the built binary in a scratch root. The stub plan-milestone workflow's first command step writes a milestone intent and a done point into RHIZA_CHANGE, and its next step runs rhiza run milestone-check, as the built-in workflow does. With reserved number 59, a plan writing milestone-59, field milestone=59 and done point m59-first succeeded. With reserved number 60, a plan writing milestone-51, field milestone=51 and done point m51-first (the old count-your-branch clash) failed at the check. Its error says the plan does not use the milestone number 60 the run reserved, and lists three problems: the wrong handle (expected milestone-60), the wrong milestone field and the done point lacking the m60- prefix. It ends with Fix: name the new milestone intent milestone-60 with field milestone=60 and prefix its done points' handles m60- (RHIZA_MILESTONE_NUMBER is 60). The plan-milestone-check script and the ./internal/api Milestone tests pass.

## Links

### Out

- verifies [milestone-check-enforces-number](../capability/milestone-check-enforces-number.md)

## History

- `01M7774MAPJ7VCPN1Q9M3RHMV4` Fix Milestone plans running at the same time pick the same milestone number and handle, so the second to merge fails CI self-check with a duplicate handle. Evidence: the agent docs plan (run 01M76EC9, record 01M76ECB) named its intent milestone-32, already held by the cleanup milestone 01M7640B, and failed pr-green-infra on PR #456. I renamed it to milestone-34 by hand (record cs 01M76V0W). The efficient build cycle plan (run 01M76STX, record 01M76SV0) had meanwhile also chosen milestone-34 and failed CI on PR #466 ('duplicate handle: node 01M76NWA... has handle milestone-34, already held by node 01M76XMNA9...'); renamed to milestone-35 by hand (record cs 01M774VK) and resume --retry. Cause: the planner takes the next number from its own branch's model, which can't see other in-flight milestone branches. Fix ideas: the daemon assigns the milestone number when it opens the record change set (counting open milestone records too), or the plan's close-out/rebase renumbers on conflict. The pre-push self-check (01M76R1E) would catch it locally only after merging main.
