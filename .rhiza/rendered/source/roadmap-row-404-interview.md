# roadmap-row-404-interview

source · `01M79SWA4CZG832Y0QVRJD3AEA`

## quotes

```json
[]
```

## raw_sha256

ddbb4541cb271b686e92560fe1ca2cc64303da75a6acf7415438223ce23fe8b4

## session

7eecefee-f9fc-4ec9-8bdc-01a1d39abb61

## summary

Autonomous plan run for bug change set 01M79SMG (roadmap rows sometimes 404). No owner was present. The planner traced the cause: loadRoadmap (internal/web/roadmap.go) reads /api/roadmap?all=true, which reaches across a root's registered projects, and links rows as /changes/<id>; under /p/<project>/ these are rewritten to that project's prefix, so another project's rows 404. Accepted change set 01M79S4M (web-roadmap-project-scoped) already scopes each project's API to that project. The planner asked how to shape this change set (depend on 01M79S4M, fold into it, or fix independently). The arbiter ruled (Decision roadmap-row-404-plan, assumed, would have escalated): plan it as a follow-on depending on 01M79S4M, with a click-every-row two-project test as acceptance, fixing any 404 that test still finds (unreadable change sets, un-prefixed links such as the change page's depends-on links). The change set serves m37-roadmap (milestone 37), whose statement says every link on the roadmap opens a page that exists.

## title

Plan interview: roadmap rows 404 (autonomous)

## Links

### In

- derived-from [roadmap-row-404-plan](../decision/roadmap-row-404-plan.md)

## History

- `01M79SMGMW0DCEBFZDHXTDTFEZ` Fix Rows on the web roadmap's backlog sometimes 404 when clicked (owner, 2026-11-13: 'Some of the items in the backlog just 404 when I click on them'). Each row links to /changes/<id> (internal/web/roadmap.go:140), but rows from another project (see the roadmap project-filter bug) or change sets the change page cannot read resolve to no page. Wanted: every roadmap row links to a page that exists in the project it belongs to (under /p/<project>/changes/<id>), and a test clicks every row of a two-project root and gets 200.
