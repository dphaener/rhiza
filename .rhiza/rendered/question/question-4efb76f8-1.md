# question-4efb76f8-1

question · `01M7AHJC72RKHKKGRWKG5783F8`

## session

4efb76f8-fe30-4c0f-81ac-01a1d5148417

## statement

The conductor's three review checks are plan (compliance + model) and code (the former quality), plus qa. The classic ship-change, the review workflow and docs-review still use the check-compliance, check-model and check-quality topics, and classic must stay unchanged. How should the plan and code checks get their role topics? / In a conductor run there is no checks step, so rhiza run findings --check and rhiza run verdict --check refuse today. How should they know which check names a conductor run accepts? / The conductor default's gate must fail while any blocker or major finding stands (not downgraded or dropped by a critic). Gate steps run commands the daemon owns. How should the gate count standing findings? / The daemon skips accepted change sets with no leaves today. Under the conductor default, the conductor plans the missing leaves. Which leafless change sets should the daemon now pick?

## Links

### In

- answers [conductor-checks-gate-pick](../decision/conductor-checks-gate-pick.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
