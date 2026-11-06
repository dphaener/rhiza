# web-overview-milestone-verified-2

evidence · `01M6QBX106H5KRC4J24RJ9JNQH` · actor [claude-code-b915ff5a](../actor/claude-code-b915ff5a.md)

## result

PASS

## statement

At e4520d10 (attempt 2): go test -trimpath ./internal/web/... passes, including TestOverviewShowsMilestoneProgress, which renders / and /fragments/milestones over two milestones with change sets in accepted and proposed states and asserts the per-state pills and counts (3, 2), the change set totals and href=/roadmap, plus the empty state. The fix for the QA finding (statement lost through lightNode) passes: TestWebMilestonesSourceNamesByStatement names two milestones by statement over a repo, the new roadmap-statement testscript and describe-golden pass, and go test ./internal/api/... passes. Hands-on, a scratch build's roadmap --all --json on this repo now gives milestone-27 and milestone-28 their statements.

## Links

### Out

- verifies [web-overview-milestone](../capability/rhiza/web-monitor/web-overview-milestone.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
