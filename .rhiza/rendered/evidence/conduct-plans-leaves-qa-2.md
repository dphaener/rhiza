# conduct-plans-leaves-qa-2

evidence · `01M7BVNWPKSBY4HR2CVFA36VMR` · actor [claude-code-348d450a](../actor/claude-code-348d450a.md)

## result

PASS

## statement

Scratch build in a scratch root: rhiza daemon --once under a conductor ship-change (first step an agent step with checks) started the leafless milestone-serving change set and skipped the leafless backlog one (serves no milestone); with a command step first, a new leafless milestone change set was skipped 'has no leaves', and rhiza run start --workflow ship-change --dry-run refused it with 'has no leaves' and next 'rhiza run start --workflow plan', while the conductor ship-change reported 'would start'; full conductor runs on leafless change sets ran to succeeded.

## Links

### Out

- verifies [conduct-plans-leaves](../capability/rhiza/orchestration/daemon/conduct-plans-leaves.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
