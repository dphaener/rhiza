# web-roadmap-project-scoped

capability · `01M79S9N7DZZ35JRWW20V549XS` · parent [web-monitor](../web-monitor.md)

## acceptance

A Go test in internal/api (beside web_root_test.go) registers two initialised projects under one root, each with its own milestone and change sets, serves WebRootHandler, GETs /p/<a>/roadmap and /p/<b>/roadmap and asserts each page lists its own project's change sets and milestone and none of the other's; it also GETs /p/<a>/api/roadmap?all=true and asserts every returned change set is project a's. go test -trimpath ./internal/api/... ./internal/web/... passes, and the existing roadmap --all testscripts still widen across the root (go test -trimpath ./cmd/rhiza -run 'TestScripts/.*roadmap.*')

## statement

Under a Rhiza root, each project's web roadmap (/p/<project>/roadmap) and change set pages list only that project's milestones and change sets, never another registered project's

## Links

### Out

- decided-by [web-api-reads-own-project](../../../decision/web-api-reads-own-project.md)

### In

- verifies [web-roadmap-project-scoped-verified](../../../evidence/web-roadmap-project-scoped-verified.md)
- verifies [web-roadmap-project-scoped-qa-1](../../../evidence/web-roadmap-project-scoped-qa-1.md)
- verifies [pr-541-fix-the-web-roadmap-page](../../../evidence/pr-541-fix-the-web-roadmap-page.md)

## History

- `01M79S4M646SATZ2HW97GXK1MH` Fix The web roadmap page served under a project (/p/rhiza/roadmap) lists change sets and milestones from every registered project, not only that project's (owner, 2026-11-13: 'The roadmap tab doesn't filter between projects properly, it shows things from both projects'). loadRoadmap (internal/web/roadmap.go:105) reads /api/roadmap?all=true through the API handler in the request context; check whether that handler, or all=true, reaches across the root's projects instead of the scoped project. Wanted: each project's roadmap shows only its own change sets and milestones, with a test serving two projects' roadmaps and asserting no cross-project rows.
