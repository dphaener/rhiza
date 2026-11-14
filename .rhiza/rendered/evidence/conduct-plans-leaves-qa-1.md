# conduct-plans-leaves-qa-1

evidence · `01M7BKZ02HXCYSNMGM3Y9KMZJ4` · actor [claude-code-73b338a9](../actor/claude-code-73b338a9.md)

## result

PASS

## statement

With the built binary: rhiza run start --workflow ship-change --dry-run on a leafless change set is refused ('has no leaves', next run start --workflow plan) under a command-step ship-change and under one whose agent step with checks is not first, and would start under the conductor ship-change; rhiza daemon --once picks nothing under the builtin ship-change, but under the conductor ship-change it starts the leafless milestone-serving change set and skips a leafless one serving a non-milestone intent.

## Links

### Out

- verifies [conduct-plans-leaves](../capability/rhiza/orchestration/daemon/conduct-plans-leaves.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
