# web-roadmap-milestone-page-qa-2

evidence · `01M7BVHJYHR4WDM2XADXFZH9J3` · actor [claude-code-e789a95e](../actor/claude-code-e789a95e.md)

## result

PASS

## statement

QA at c1769e866, a scratch-build daemon (paused) on a two-project root: /p/rhiza/roadmap showed 30 open milestones; each one's kicker state and bar aria-label (landed, accepted, proposed and withdrawn counts) matched /api/milestones for all 30 (checked by script). Done points sat in closed folds, groups ran Accepted waiting, Proposed, Landed, Drafts, landed rows were titled from the title field ('Builder sessions' actor'), the backlog and the umbrella came after the milestones, '6 finished milestones' was one in-page fold with no links, no finished milestone was among the open ones and every link stayed under /p/rhiza. /p/other/roadmap said Nothing is planned. Building now was not exercised (no runs in the clone). The page took 50-91s per load at host load ~17.

## Links

### Out

- verifies [web-roadmap-milestone-page](../capability/rhiza/web-monitor/web-roadmap-milestone-page.md)

## History

- `01M79NMNGDCNJPAWRP791PRF7G` Web roadmap and change set page: milestones, their done points and change sets with short titles, state and progress, scoped to the project; a change set page with its leaves, acceptance, Evidence, decisions, PR and runs
