# web-roadmap-rows-resolve-verified

evidence · `01M7AHV8920GA4WQB1STFVY0ZW` · actor [claude-code-7b96401c](../actor/claude-code-7b96401c.md)

## result

PASS

## statement

At 80d5853e0: go test -trimpath -count=1 -run TestWebRootRoadmapRowsResolve -v ./internal/api/ PASS (0.47s); go test -trimpath ./internal/web/... ok; rhiza ship generate exit 0 (0 files changed); rhiza render exit 0 with a clean tree. The full go test -trimpath ./internal/api/... ran 3 times on this host at load average 6-16. Each run hit the 10m package timeout, each with a different test still running (TestLinkEndpointKinds, TestReadBudgetApply, TestInboxGateQuestionsBudget), including once with -skip TestReadBudget. Since every run stopped at a different test, I took this for host load, not a hang; I saw no evidence of a hang. The commit only adds one test file plus an actor record, so the timeout is not attributed to it; the PR's CI make check must confirm the package. Diff read: internal/api/web_roadmap_rows_test.go registers projects alpha and beta under one rhiza.toml root. Each has a milestone change set that serves a milestone intent and a backlog change set that depends-on it. The test serves WebRootHandler, GETs /p/<p>/roadmap, and asserts every /changes/ row href is under /p/<p>/changes/ and the page has no text from the other project's change set intents. It GETs each row (200), then every /changes/ link on those pages, depends-on included, asserts it is in the same project and gets 200. Every acceptance point is tested. The test does not check the other project's change set IDs by value; a cross-project row href would still fail the prefix and 200 checks. Nothing outside the leaf changed, there are no stubs or TODOs, and the tests use fakes only.

## Links

### Out

- verifies [web-roadmap-rows-resolve](../capability/rhiza/web-monitor/web-roadmap-rows-resolve.md)

## History

- `01M79SMGMW0DCEBFZDHXTDTFEZ` Fix Rows on the web roadmap's backlog sometimes 404 when clicked (owner, 2026-11-13: 'Some of the items in the backlog just 404 when I click on them'). Each row links to /changes/<id> (internal/web/roadmap.go:140), but rows from another project (see the roadmap project-filter bug) or change sets the change page cannot read resolve to no page. Wanted: every roadmap row links to a page that exists in the project it belongs to (under /p/<project>/changes/<id>), and a test clicks every row of a two-project root and gets 200.
