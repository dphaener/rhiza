# web-overview-milestone-verified

evidence · `01M6Q8H2K3XRYXE66J8MB8TKBF` · actor [claude-code-0c886717](../actor/claude-code-0c886717.md)

## result

PASS

## statement

At b0e3fa7f, TestOverviewShowsMilestoneProgress in internal/web renders / (and /fragments/milestones) through HandlerWithMilestones over a fixture milestone with 5 change sets (3 accepted, 2 proposed) plus a one-change-set milestone. It asserts the rz-pill--accepted and rz-pill--proposed pills with counts 3 and 2, the totals ('5 change sets', '1 change set'), '2 open' and href="/roadmap", plus the empty state. go test -trimpath ./internal/web/ passed (overview, layout and partials goldens). go test ./internal/api -run 'Roadmap|Web|Plan|Daemon' passed, and so did cmd/rhiza TestScripts/describe-golden, which carries the new RoadmapIntent.statement. From the diff: the daemon feeds the panel from runRoadmap (All) milestone groups, named by statement with no IDs.

## Links

### Out

- verifies [web-overview-milestone](../capability/rhiza/web-monitor/web-overview-milestone.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
