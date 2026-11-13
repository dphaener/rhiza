# web-roadmap-project-scoped-verified

evidence · `01M7A6NFMP4B6QKX63ZY2XFVYR` · actor [claude-code-8bbb4dd6](../actor/claude-code-8bbb4dd6.md)

## result

PASS

## statement

Verified at 003b4ffad: go test -trimpath ./internal/api/... ./internal/web/... exit 0; go test -trimpath ./cmd/rhiza -run 'TestScripts/.*roadmap.*' exit 0 (roadmap --all scripts still widen across the root); rhiza ship generate and rhiza render exit 0 with no changes to tracked files. New TestWebRootRoadmapIsProjectScoped (internal/api/web_roadmap_scope_test.go) registers alpha and beta under one root with their own milestone and change sets, serves WebRootHandler, and checks that /p/<x>/roadmap lists only its own milestone and change sets and that /p/<x>/api/roadmap?all=true returns only its own 2 change sets. Reverting the web_api.go fix (ReadOnlyHandler forces RoadmapRequest.local) makes the test fail with beta's API returning alpha's change sets (4, want 2). The diff touches only web_api.go, the new test and run records; it has no stubs or TODOs.

## Links

### Out

- verifies [web-roadmap-project-scoped](../capability/rhiza/web-monitor/web-roadmap-project-scoped.md)

## History

- `01M79S4M646SATZ2HW97GXK1MH` Fix The web roadmap page served under a project (/p/rhiza/roadmap) lists change sets and milestones from every registered project, not only that project's (owner, 2026-11-13: 'The roadmap tab doesn't filter between projects properly, it shows things from both projects'). loadRoadmap (internal/web/roadmap.go:105) reads /api/roadmap?all=true through the API handler in the request context; check whether that handler, or all=true, reaches across the root's projects instead of the scoped project. Wanted: each project's roadmap shows only its own change sets and milestones, with a test serving two projects' roadmaps and asserting no cross-project rows.
