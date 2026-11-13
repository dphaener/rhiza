# web-api-reads-own-project

decision · `01M79S9N480458ETJ4BRJPKPHF`

## statement

A project's web API (the read-only handler a project's web UI reads through) reads only that project: roadmap with all=true returns every change set of the project it serves and never widens to the root's other projects; the CLI's rhiza roadmap --all keeps widening across the root

## wording

Scope the web API to its project, chosen by the planner in an autonomous run

## Links

### Out

- derived-from [plan-web-roadmap-project-scope](../source/plan-web-roadmap-project-scope.md)

### In

- decided-by [web-roadmap-project-scoped](../capability/rhiza/web-monitor/web-roadmap-project-scoped.md)

## History

- `01M79S4M646SATZ2HW97GXK1MH` Fix The web roadmap page served under a project (/p/rhiza/roadmap) lists change sets and milestones from every registered project, not only that project's (owner, 2026-11-13: 'The roadmap tab doesn't filter between projects properly, it shows things from both projects'). loadRoadmap (internal/web/roadmap.go:105) reads /api/roadmap?all=true through the API handler in the request context; check whether that handler, or all=true, reaches across the root's projects instead of the scoped project. Wanted: each project's roadmap shows only its own change sets and milestones, with a test serving two projects' roadmaps and asserting no cross-project rows.
