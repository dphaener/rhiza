# web-links-resolve-verified

evidence · `01M7BGAAM3N56H5CTAE6GAB37G` · actor [claude-code-71b9458a](../actor/claude-code-71b9458a.md)

## result

PASS

## statement

At 3e2f3f6b2: go test -trimpath ./internal/web -run TestLinksResolve -count=1 PASS (exit 0); go test -trimpath ./internal/web/... ok; rhiza ship generate and rhiza render exit 0 with no tracked changes. internal/web/links_test.go builds a two-project root (a, b) via RootWithProjects, each with a landed change set (intent, capability, decision) and a run, and crawls /p/<name>/ Home, /roadmap, /changes/<cs>, /runs/<run> and /nodes/links-cap, following every href/src/data-rz-fragment link: each in-project link must return 200, any link leaving the project fails unless it is a shared root path, and at least 8 pages must be crawled. Fakes only (no network, Rex or gh); diff touches only the test file plus actor log/rendered records; no stubs or TODOs.

## Links

### Out

- verifies [web-links-resolve](../capability/rhiza/web-monitor/web-links-resolve.md)

## History

- `01M79NN43PGK6HDXTRY33VNC1E` Web model pages: every capability, decision and source link opens a designed page built from the model, with its statement, its why (decisions, sources, change sets) and its history
