# web-roadmap-project-scoped-qa-1

evidence · `01M7A73WH09GD8DA69AAJV867X` · actor [claude-code-939067d5](../actor/claude-code-939067d5.md)

## result

PASS

## statement

QA at fe1ef0b8d: built ./cmd/rhiza to a scratch path, made a scratch root with projects alpha and beta (each one milestone and two proposed change sets, beta's second depending on its first) and ran the built daemon, paused with max-runs 0, serving the web UI on 127.0.0.1. /p/alpha/roadmap and /p/beta/roadmap each list only their own two change sets. /api/roadmap?all=true (home) and /p/<x>/api/roadmap?all=true return only that project's milestone and 2 change sets, and so do ?project=alpha and ?all=true&project=alpha on beta's API. /p/beta/changes/<beta two> shows Serves beta milestone and Depends on beta change one, and /p/beta/changes/<an alpha change set> answers 404. Failure cases: /p/gamma/roadmap 404, ?all=maybe 400 usage error, POST 405. The CLI's rhiza roadmap --all still lists both projects. go test -trimpath ./internal/api -run TestWebRoot, ./internal/web/... and ./cmd/rhiza -run 'TestScripts/.*roadmap.*' pass.

## Links

### Out

- verifies [web-roadmap-project-scoped](../capability/rhiza/web-monitor/web-roadmap-project-scoped.md)

## History

- `01M79S4M646SATZ2HW97GXK1MH` Fix The web roadmap page served under a project (/p/rhiza/roadmap) lists change sets and milestones from every registered project, not only that project's (owner, 2026-11-13: 'The roadmap tab doesn't filter between projects properly, it shows things from both projects'). loadRoadmap (internal/web/roadmap.go:105) reads /api/roadmap?all=true through the API handler in the request context; check whether that handler, or all=true, reaches across the root's projects instead of the scoped project. Wanted: each project's roadmap shows only its own change sets and milestones, with a test serving two projects' roadmaps and asserting no cross-project rows.
