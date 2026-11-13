# web-roadmap-rows-resolve

capability · `01M79SWNGATEE7JHYB7Q9X2CBK` · parent [web-monitor](../web-monitor.md)

## acceptance

A Go test in internal/api (beside web_root_test.go) registers two initialised projects under one root, each with a milestone change set and a backlog change set (one depending on another), serves WebRootHandler, GETs /p/<a>/roadmap and /p/<b>/roadmap, extracts every row href, asserts each is under that project's /p/<project>/changes/ prefix and names none of the other project's change sets, GETs each and gets 200, then GETs every depends-on link on those change pages and gets 200. go test -trimpath ./internal/api/... ./internal/web/... passes

## statement

Under a Rhiza root, every row on each project's web roadmap (/p/<project>/roadmap: milestone, umbrella and backlog groups) and every depends-on link on its change set pages opens that project's change set page (/p/<project>/changes/<id>), never a 404

## Links

### Out

- decided-by [roadmap-row-404-plan](../../../decision/roadmap-row-404-plan.md)

### In

- verifies [web-roadmap-rows-resolve-verified](../../../evidence/web-roadmap-rows-resolve-verified.md)
- verifies [web-roadmap-rows-resolve-qa-1](../../../evidence/web-roadmap-rows-resolve-qa-1.md)
- verifies [pr-549-fix-rows-on-the-web-roadmap-s](../../../evidence/pr-549-fix-rows-on-the-web-roadmap-s.md)

## History

- `01M79SMGMW0DCEBFZDHXTDTFEZ` Fix Rows on the web roadmap's backlog sometimes 404 when clicked (owner, 2026-11-13: 'Some of the items in the backlog just 404 when I click on them'). Each row links to /changes/<id> (internal/web/roadmap.go:140), but rows from another project (see the roadmap project-filter bug) or change sets the change page cannot read resolve to no page. Wanted: every roadmap row links to a page that exists in the project it belongs to (under /p/<project>/changes/<id>), and a test clicks every row of a two-project root and gets 200.
