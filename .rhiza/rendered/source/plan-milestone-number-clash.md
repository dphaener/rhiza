# plan-milestone-number-clash

source · `01M777AW50YHP0C4A0D7YQET8E`

## quotes

```json
[]
```

## raw_sha256

1162754a2e64c726b57ee48f9d8b1a80c6b069165f2c2c98e9d6e865f1e4726d

## session

d85dd023-7452-4f12-87c3-01a1ce728739

## summary

Autonomous plan run of bug change set 01M7774: two milestone plans running at once both took the next milestone number from their own branch's model and clashed in CI (milestone-32 on PR #456, milestone-34 on PR #466), fixed by hand. The planner asked the arbiter three questions in one: which intent the fix serves, how to stop the clash, and whether run.milestone-check should enforce the number. The arbiter ruled (Decision milestone-number-reserved-at-run-start): serve milestone-25; run start of a milestone plan reserves N = 1 + the highest of landed milestone numbers and every number reserved in any run journal in the shared run store (ended and failed runs included), stores it in the journal, shows it in run status and hands it to the planner as RHIZA_MILESTONE_NUMBER; the plan guide tells the planner to use handle milestone-N, field milestone=N and done-point prefix mN-; a run planning an existing milestone reserves nothing; no renumbering at close-out; run.milestone-check fails a plan whose new milestone intent does not use the reserved number, naming the intent and expected number. Three leaves follow: the reservation, the planner's handoff and guide, and the check.

## title

Plan of the concurrent milestone number bug, with the arbiter

## Links

### In

- derived-from [milestone-number-reserved-at-run-start](../decision/milestone-number-reserved-at-run-start.md)

## History

- `01M7774MAPJ7VCPN1Q9M3RHMV4` Fix Milestone plans running at the same time pick the same milestone number and handle, so the second to merge fails CI self-check with a duplicate handle. Evidence: the agent docs plan (run 01M76EC9, record 01M76ECB) named its intent milestone-32, already held by the cleanup milestone 01M7640B, and failed pr-green-infra on PR #456. I renamed it to milestone-34 by hand (record cs 01M76V0W). The efficient build cycle plan (run 01M76STX, record 01M76SV0) had meanwhile also chosen milestone-34 and failed CI on PR #466 ('duplicate handle: node 01M76NWA... has handle milestone-34, already held by node 01M76XMNA9...'); renamed to milestone-35 by hand (record cs 01M774VK) and resume --retry. Cause: the planner takes the next number from its own branch's model, which can't see other in-flight milestone branches. Fix ideas: the daemon assigns the milestone number when it opens the record change set (counting open milestone records too), or the plan's close-out/rebase renumbers on conflict. The pre-push self-check (01M76R1E) would catch it locally only after merging main.
