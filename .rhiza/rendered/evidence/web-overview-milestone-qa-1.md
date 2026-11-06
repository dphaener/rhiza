# web-overview-milestone-qa-1

evidence · `01M6Q95MV2RAA1X98FGN8REX7N` · actor [claude-code-885ce195](../actor/claude-code-885ce195.md)

## result

FAIL

## statement

Same live daemon over this repo's real model: the Milestones panel showed '2 open' with Roadmap link /roadmap and per-milestone counts by state (proposed 2 of 2; accepted 3 of 3), but both rows were named just 'Milestone': rhiza roadmap --all --json gives milestone-27 and milestone-28 no statement, though rhiza show milestone-28 has one. planRead keeps the intents through lightNode, which drops Fields, so RoadmapIntent.Statement is always empty from a real repo.

## Links

### Out

- verifies [web-overview-milestone](../capability/rhiza/web-monitor/web-overview-milestone.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
