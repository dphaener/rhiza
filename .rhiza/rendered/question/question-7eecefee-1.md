# question-7eecefee-1

question · `01M79SSW417MJXAT0Y0NXGGR1J`

## session

7eecefee-f9fc-4ec9-8bdc-01a1d39abb61

## statement

Roadmap rows 404 because /api/roadmap?all=true reaches across the root's projects, so a project's roadmap shows another project's change sets and their rewritten /p/<this>/changes/<id> links aren't found. Accepted change set 01M79S4M (web-roadmap-project-scoped) already scopes each project's API to that project. How should this bug change set 01M79SMG be planned?

## Links

### In

- answers [roadmap-row-404-plan](../decision/roadmap-row-404-plan.md)

## History

- `01M79SMGMW0DCEBFZDHXTDTFEZ` Fix Rows on the web roadmap's backlog sometimes 404 when clicked (owner, 2026-11-13: 'Some of the items in the backlog just 404 when I click on them'). Each row links to /changes/<id> (internal/web/roadmap.go:140), but rows from another project (see the roadmap project-filter bug) or change sets the change page cannot read resolve to no page. Wanted: every roadmap row links to a page that exists in the project it belongs to (under /p/<project>/changes/<id>), and a test clicks every row of a two-project root and gets 200.
