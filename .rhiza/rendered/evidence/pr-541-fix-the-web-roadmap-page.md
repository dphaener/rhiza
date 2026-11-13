# pr-541-fix-the-web-roadmap-page

evidence · `01M7A8PCCQNSNBT4ZS8STBS9AF` · actor [dphaener](../actor/dphaener.md)

## head

896ff186072836830f297a92d896e5f0b0aa1400

## pr

https://github.com/dphaener/rhiza/pull/541

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #541 (https://github.com/dphaener/rhiza/pull/541) passed CI at 896ff18

## status

verified

## Links

### Out

- verifies [web-roadmap-project-scoped](../capability/rhiza/web-monitor/web-roadmap-project-scoped.md)

## History

- `01M79S4M646SATZ2HW97GXK1MH` Fix The web roadmap page served under a project (/p/rhiza/roadmap) lists change sets and milestones from every registered project, not only that project's (owner, 2026-11-13: 'The roadmap tab doesn't filter between projects properly, it shows things from both projects'). loadRoadmap (internal/web/roadmap.go:105) reads /api/roadmap?all=true through the API handler in the request context; check whether that handler, or all=true, reaches across the root's projects instead of the scoped project. Wanted: each project's roadmap shows only its own change sets and milestones, with a test serving two projects' roadmaps and asserting no cross-project rows.
