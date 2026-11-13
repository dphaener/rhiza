# plan-finding-kinds

capability · `01M79BH55KDAWPGWR01WXZ38N3` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/run/... passes with tests that run.Kinds ends with uncovered-intent then oversized-leaf and KindSeverity returns major for both; go test -trimpath ./internal/api -run 'RunFindings|ReviewCheck' passes with a round holding a confirmed uncovered-intent finding failing review-check before round 3, and (failure case) a finding of an unknown kind such as no-failure-case still refused with an error listing the seven kinds; go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden passes after regenerating it with -update.

## statement

`rhiza run findings` accepts two more plan-review kinds, `uncovered-intent` (a clause of the change set's intent that no leaf serves) and `oversized-leaf` (a leaf with more than one concern, or too large for one fresh session), listed after the five existing kinds in run.Kinds, and `rhiza run review-check` treats both as major, as its summary and the findings schema say.

## Links

### Out

- decided-by [plan-role-topics-shape](../../../../../decision/plan-role-topics-shape.md)
- derived-from [plan-role-topics-interview](../../../../../source/plan-role-topics-interview.md)

### In

- depends-on [shared-severities-block](shared-severities-block.md)
- depends-on [plan-reviewer-topic](plan-reviewer-topic.md)
- depends-on [plan-reviser-topic](plan-reviser-topic.md)
- verifies [plan-finding-kinds-verified](../../../../../evidence/plan-finding-kinds-verified.md)
- verifies [plan-finding-kinds-qa-1](../../../../../evidence/plan-finding-kinds-qa-1.md)
- verifies [plan-finding-kinds-verified-5](../../../../../evidence/plan-finding-kinds-verified-5.md)
- verifies [pr-533-plan-role-topics](../../../../../evidence/pr-533-plan-role-topics.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
