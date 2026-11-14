# conduct-skill-plan-2026-11-13

source · `01M7AHT653JXZV2Y1YH0CBNQWC`

## quotes

```json
[]
```

## raw_sha256

bec1f960ed44448df70322365293d8a686fe795d45789cd85df67b51d5e4567b

## session

4efb76f8-fe30-4c0f-81ac-01a1d5148417

## summary

Autonomous plan run for change set 01M79KGZ (A rhiza-conduct skill conducts), milestone 36. There was no owner. The research found four things: the no-leaves skip is in internal/api daemonCandidates, not internal/daemon; rhiza run findings --check and run verdict --check accept only a checks step's declared checks; the critic-on-blocker-or-major rule already exists as run.Round.Standing; and the classic ship-change, review and docs-review still load check-compliance, check-model and check-quality. Four questions went to the arbiter: the role topics for the plan and code checks, how a conductor run knows its check names, how the gate counts standing findings, and which leafless change sets the daemon picks. The arbiter answered them all in Decision conductor-checks-gate-pick: new check-plan and check-code topics beside the old ones; the conduct step declares its checks the way a checks step does, and findings and verdict validate against them; the gate runs rhiza run review-check, extended to conductor runs and counted from the journal; a leafless, milestone-serving change set is picked only when its dispatch resolves to the conductor default. The plan was rewritten into five leaves, and the change set depends on the spawn and contract change sets.

## title

Planning the conduct skill change set, autonomous, with the arbiter

## Links

### In

- derived-from [conductor-checks-gate-pick](../decision/conductor-checks-gate-pick.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
