# roadmap-row-404-plan

decision · `01M79SVQVPME0SNAWE3XWG7HKN`

## reasoning

The question's diagnosis matches the code: loadRoadmap (internal/web/roadmap.go:105) reads /api/roadmap?all=true and links rows as /changes/<id> (roadmap.go:140), so cross-project rows get rewritten to /p/<this>/changes/<id> and 404. Accepted 01M79S4M already owns that root cause (scope each project's roadmap API to that project, with a no-cross-project-rows test). Duplicating the scoping here would make two change sets edit the same function. But this bug's intent asks for more than 01M79S4M promises: every row's link resolves (200) in a two-project root, which also covers change sets the change page cannot read and other un-prefixed links (change.go:178 depends-on, inbox.go:44). So the bug stays open as a dependent change set whose acceptance is the click-every-row 200 test, fixing whatever residual 404 it exposes.

## statement

Plan 01M79SMG as a follow-on to accepted 01M79S4M (web-roadmap-project-scoped), not a second fix of the same cause: it depends on 01M79S4M and does not re-scope /api/roadmap itself. Its leaf's acceptance is the regression the intent asks for: a test serves a root with two registered projects, renders each project's roadmap, follows every row link (and the change page's depends-on links) under /p/<project>/changes/<id>, and gets 200 for each, with no row from the other project. Any 404 that test still finds once 01M79S4M has landed (a change set the change page cannot read, a link not rewritten under the project prefix) is fixed in this change set's leaf.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-7eecefee-1](../question/question-7eecefee-1.md)
- derived-from [roadmap-row-404-interview](../source/roadmap-row-404-interview.md)

### In

- decided-by [web-roadmap-rows-resolve](../capability/rhiza/web-monitor/web-roadmap-rows-resolve.md)

## History

- `01M79SMGMW0DCEBFZDHXTDTFEZ` Fix Rows on the web roadmap's backlog sometimes 404 when clicked (owner, 2026-11-13: 'Some of the items in the backlog just 404 when I click on them'). Each row links to /changes/<id> (internal/web/roadmap.go:140), but rows from another project (see the roadmap project-filter bug) or change sets the change page cannot read resolve to no page. Wanted: every roadmap row links to a page that exists in the project it belongs to (under /p/<project>/changes/<id>), and a test clicks every row of a two-project root and gets 200.
